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
)

// ExecNotification is used to execute an RPC method regarding NETCONF notification
func (s *Session) ExecNotification(operation interface{}) (*message.Notification, error) {
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

	fmt.Println("\nReceiving Notification's answer")
	rawXML, err := s.Transport.Receive()
	if err != nil {
		return nil, err
	}

	reply, err := message.NewNotification(rawXML)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
