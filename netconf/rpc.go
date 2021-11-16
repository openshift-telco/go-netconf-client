// Go NETCONF Client
//
// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

import (
	"encoding/xml"
	"fmt"
	"github.com/adetalhouet/go-netconf/netconf/message"
	"strings"
)

// ExecRPC is used to execute an RPC method
func (s *Session) ExecRPC(operation interface{}) (*message.RPCReply, error) {
	request, err := xml.Marshal(operation)
	if err != nil {
		return nil, err
	}

	header := []byte(xml.Header)
	request = append(header, request...)

	fmt.Println(fmt.Sprintf("\n\nSending RPC"))
	err = s.Transport.Send(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("\nReceiving RPC's answer")
	rawXML, err := s.Transport.Receive()
	if err != nil {
		return nil, err
	}

	// FIXME better handle different messages:  notification and rpc-reply.
	// For now, we discard notification, as they *should* be handled separately.
	// In a case of establish-subscription, sometimes the device sends the notification
	// before sending the rpc-reply.
	var rawReply = string(rawXML)
	if strings.Contains(rawReply, "<rpc-reply") {
		return message.NewRPCReply(rawXML)
	}

	return nil, fmt.Errorf("discarding received message: %s", rawReply)
}
