// Go NETCONF Client
//
// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

const (
	// msgSeparator is used to separate sent messages via NETCONF
	msgSeparator    = "]]>]]>"
	msgSeparatorV11 = "\n##\n"
)

// Transport interface defines what characteristics make up a NETCONF transport
// layer object.
type Transport interface {
	Send([]byte) error
	Receive() ([]byte, error)
	Close() error
	SetVersion(version string)
}

type transportBasicIO struct {
	io.ReadWriteCloser
	//new add
	version string
}

func (t *transportBasicIO) SetVersion(version string) {
	t.version = version
}

// Send a well formatted NETCONF rpc message as a slice of bytes adding on the
// necessary framing messages.
func (t *transportBasicIO) Send(data []byte) error {
	var separator []byte
	var dataInfo []byte
	if t.version == "v1.1" {
		separator = append(separator, []byte(msgSeparatorV11)...)
	} else {
		separator = append(separator, []byte(msgSeparator)...)
	}

	if t.version == "v1.1" {
		header := fmt.Sprintf("\n#%d\n", len(string(data)))
		dataInfo = append(dataInfo, header...)
	}
	dataInfo = append(dataInfo, data...)
	dataInfo = append(dataInfo, separator...)
	_, err := t.Write(dataInfo)

	println(string(dataInfo))

	return err
}

func (t *transportBasicIO) Receive() ([]byte, error) {
	var separator []byte
	if t.version == "v1.1" {
		separator = append(separator, []byte(msgSeparatorV11)...)
		// NOTES: This is not clever at all
		// you are reading the O-RU response content once with WaitForBytes, and then you read it again to get rid of
		// the #<chunk-size> pieces. Using Chunked would be enough, but if you pass in the t.ReadWriteCloser to the
		// splitChunked function it gets stuck when doing the Read of the last piece of the NETCONF message.
		// This will need to be addressed in the future.
		// Also, splitChunked function comes from the andaru/netconf library, with just a slight modification to make it
		// work.
		b, err := t.WaitForBytes(separator)
		if err != nil {
			return nil, err
		}
		return t.Chunked(b)
	}
	separator = append(separator, []byte(msgSeparator)...)
	return t.WaitForBytes(separator)
}

func (t *transportBasicIO) Writeln(b []byte) (int, error) {
	_, err := t.Write(b)
	if err != nil {
		return 0, err
	}
	_, err = t.Write([]byte("\n"))
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// SplitChunked returns a bufio.SplitFunc suitable for decoding
// "chunked framing" NETCONF transport streams.
//
// endOfMessage will be called at the end of each NETCONF message,
// and must not be nil.
//
// It must only be used with bufio.Scanner who have a buffer of
// at least 16 bytes (rarely is this a concern).
func SplitChunked(endOfMessage func()) bufio.SplitFunc {
	type stateT int
	const (
		headerStart stateT = iota
		headerSize
		data
		endOfChunks
	)
	var state stateT
	var cs, dataleft int

	return func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		for cur := b[advance:]; err == nil && advance < len(b); cur = b[advance:] {
			if len(cur) < 4 && !atEOF {
				return
			}
			switch state {
			case headerStart:
				switch {
				case bytes.HasPrefix(cur, []byte("\n#")):
					if len(cur) < 4 {
						err = ErrBadChunk
						return
					}
					switch r := cur[2]; {
					case r == '#':
						advance += 3
						state = endOfChunks
					case r >= '1' && r <= '9':
						advance += 2
						state = headerSize
					default:
						err = ErrBadChunk
					}
				default:
					err = ErrBadChunk
				}
			case headerSize:
				switch idx := bytes.IndexByte(cur, '\n'); {
				case idx < 1, idx > 10:
					if len(cur) < 11 && !atEOF {
						return
					}
					err = ErrBadChunk
				default:
					csize := cur[:idx]
					if csizeVal, csizeErr := strconv.ParseUint(string(csize), 10, 31); csizeErr != nil {
						err = ErrBadChunk
					} else {
						advance += idx + 1
						dataleft = int(csizeVal)
						state = data
					}
				}
			case data:
				var rsize int
				if rsize = len(cur); dataleft < rsize {
					rsize = dataleft
				}
				token = append(token, cur[:rsize]...)
				advance += rsize
				if dataleft -= rsize; dataleft < 1 {
					state = headerStart
					cs++
				}
				if rsize > 0 {
					return
				}
			case endOfChunks:
				switch r := cur[0]; {
				case r == '\n' && cs > 0:
					advance++
					state = headerStart
					if endOfMessage != nil {
						endOfMessage()
					}
				default:
					err = ErrBadChunk
				}
			}
		}
		if atEOF && dataleft > 0 {
			//state = headerStart
			err = io.ErrUnexpectedEOF
		}
		return
	}
}

// ErrBadChunk indicates a chunked framing protocol error occurred
var ErrBadChunk = errors.New("bad chunk")

func (t *transportBasicIO) Chunked(b []byte) ([]byte, error) {
	rdr := bytes.NewReader(b)
	scanner := bufio.NewScanner(rdr)
	bsize := 16
	scanner.Buffer(make([]byte, bsize), bsize*2)

	scanner.Split(SplitChunked(nil))
	var got []byte
	for scanner.Scan() {
		got = append(got, scanner.Bytes()...)
	}
	return got, nil
}

func (t *transportBasicIO) WaitForFunc(f func([]byte) (int, error)) ([]byte, error) {
	var out bytes.Buffer
	buf := make([]byte, 8192)

	pos := 0
	for {
		n, err := t.Read(buf[pos : pos+(len(buf)/2)])
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		if n > 0 {
			end, err := f(buf[0 : pos+n])
			if err != nil {
				return nil, err
			}

			if end > -1 {
				out.Write(buf[0:end])
				return out.Bytes(), nil
			}

			if pos > 0 {
				out.Write(buf[0:pos])
				copy(buf, buf[pos:pos+n])
			}

			pos = n
		}
	}

	return nil, fmt.Errorf("WaitForFunc failed")
}

func (t *transportBasicIO) WaitForBytes(b []byte) ([]byte, error) {
	return t.WaitForFunc(
		func(buf []byte) (int, error) {
			return bytes.Index(buf, b), nil
		},
	)
}

func (t *transportBasicIO) WaitForString(s string) (string, error) {
	out, err := t.WaitForBytes([]byte(s))
	if out != nil {
		return string(out), err
	}
	return "", err
}

func (t *transportBasicIO) WaitForRegexp(re *regexp.Regexp) ([]byte, [][]byte, error) {
	var matches [][]byte
	out, err := t.WaitForFunc(
		func(buf []byte) (int, error) {
			loc := re.FindSubmatchIndex(buf)
			if loc != nil {
				for i := 2; i < len(loc); i += 2 {
					matches = append(matches, buf[loc[i]:loc[i+1]])
				}
				return loc[1], nil
			}
			return -1, nil
		},
	)
	return out, matches, err
}

// ReadWriteCloser represents a combined IO Reader and WriteCloser
type ReadWriteCloser struct {
	io.Reader
	io.WriteCloser
}

// NewReadWriteCloser creates a new combined IO Reader and WriteCloser from the provided objects
func NewReadWriteCloser(r io.Reader, w io.WriteCloser) *ReadWriteCloser {
	return &ReadWriteCloser{r, w}
}
