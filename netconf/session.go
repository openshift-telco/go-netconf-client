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

	return s
}

// SendHello send the initial message through NETCONF to advertise supported capability.
func (s *Session) SendHello(hello *message.Hello) error {
	val, err := xml.Marshal(hello)
	if err != nil {
		return err
	}

	header := []byte(xml.Header)
	val = append(header, val...)
	err = s.Transport.Send(val)

	// Set Transport version after sending hello-message,
	// so the hello-message is sent using netconf:1.0 framing
	s.Transport.SetVersion("v1.0")
	for _, capability := range s.Capabilities {
		if strings.Contains(capability, message.NetconfVersion11) {
			s.Transport.SetVersion("v1.1")
			break
		}
	}

	return err
}

// ReceiveHello is the first message received when connecting to a NETCONF server.
// It provides the supported capabilities of the server.
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
