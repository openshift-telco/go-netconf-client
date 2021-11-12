// Go NETCONF Client
//
// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

import (
	"encoding/xml"
	"github.com/adetalhouet/go-netconf/netconf/message"
	"strings"
)

// DefaultCapabilities sets the default capabilities of the client library
var DefaultCapabilities = []string{
	message.NetconfVersion10,
	message.NetconfVersion11,
	"urn:ietf:params:xml:ns:yang:ietf-event-notifications",
	"urn:ietf:params:xml:ns:yang:ietf-yang-push",
}

// Session represents a NETCONF sessions with a remote NETCONF server
type Session struct {
	Transport    Transport
	SessionID    int
	Capabilities []string
}

// NewSession creates a new NETCONF session using the provided transport layer.
func NewSession(t Transport) *Session {
	s := new(Session)
	s.Transport = t

	// Receive server Hello message
	serverHello, _ := s.ReceiveHello()
	s.SessionID = serverHello.SessionID
	s.Capabilities = serverHello.Capabilities

	// Set Transport version
	t.SetVersion("v1.0")
	for _, capability := range s.Capabilities {
		if strings.Contains(capability, message.NetconfVersion11) {
			t.SetVersion("v1.1")
			break
		}
	}

	return s
}

func (s *Session) SendHello(hello *message.Hello) error {
	val, err := xml.Marshal(hello)
	if err != nil {
		return err
	}

	header := []byte(xml.Header)
	val = append(header, val...)
	err = s.Transport.Send(val)
	return err
}

func (s *Session) ReceiveHello() (*message.Hello, error) {
	hello := new(message.Hello)

	val, err := s.Transport.Receive()
	if err != nil {
		return hello, err
	}

	err = xml.Unmarshal(val, hello)
	return hello, err
}

// Close is used to close and end a session
func (s *Session) Close() error {
	return s.Transport.Close()
}
