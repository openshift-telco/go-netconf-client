// Go NETCONF Client
//
// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

import (
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
	serverHello, _ := t.ReceiveHello()
	s.SessionID = serverHello.SessionID
	s.Capabilities = serverHello.Capabilities

	// Send our hello using default capabilities.
	err := t.SendHello(&message.Hello{Capabilities: DefaultCapabilities})
	if err != nil {
		return nil
	}

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

// Close is used to close and end a session
func (s *Session) Close() error {
	return s.Transport.Close()
}