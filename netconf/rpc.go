// Go NETCONF Client
//
// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/adetalhouet/go-netconf/netconf/message"
	"strings"
)

// ExecRPC is used to execute an RPC method
func (s *Session) ExecRPC(operation interface{}) (interface{}, error) {
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

	var rawReply = string(rawXML)
	if strings.Contains(rawReply, "<rpc-reply") {
		return message.NewRPCReply(rawXML)
	} else if strings.Contains(rawReply, "<notification") {
		return message.NewNotification(rawXML)
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown received message: %s", rawReply))
	}
}
