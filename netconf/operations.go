/*
Copyright 2021. Alexis de Talhouët

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package netconf

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"

	"github.com/openshift-telco/go-netconf-client/netconf/message"
)

// CreateNotificationStream is a convenient method to create a notification stream registration.
// TODO limitation - for now, we can only register one stream per session, because when a notification is received
// there is no way to attribute it to a specific stream
func (session *Session) CreateNotificationStream(
	timeout int32, stopTime string, startTime string, filter string, stream string, callback Callback,
) error {
	if session.IsNotificationStreamCreated {
		return fmt.Errorf(
			"there is already an active notification stream subscription. " +
				"A session can only support one notification stream at the time",
		)
	}
	session.Listener.Register(message.NetconfNotificationStreamHandler, callback)
	sub := message.NewCreateSubscription(stopTime, startTime, stream, filter)
	rpc, err := session.SyncRPC(sub, timeout)
	if err != nil {
		errMsg := "fail to create notification stream"
		if rpc != nil && len(rpc.Errors) != 0 {
			errMsg += fmt.Sprintf(" with errors: %s", rpc.Errors)
		}
		return fmt.Errorf("%s: %w", errMsg, err)
	}
	session.IsNotificationStreamCreated = true
	return nil
}

// AsyncRPC is used to send an RPC method and receive the response asynchronously.
func (session *Session) AsyncRPC(operation message.RPCMethod, callback Callback) error {

	// get XML payload
	request, err := marshall(operation)
	if err != nil {
		return err
	}

	// register the listener for the message
	session.Listener.Register(operation.GetMessageID(), callback)

	session.logger.Info("Sending RPC")
	err = session.Transport.Send(request)
	if err != nil {
		return err
	}

	return nil
}

// SyncRPC is used to execute an RPC method and receive the response synchronously
func (session *Session) SyncRPC(operation message.RPCMethod, timeout int32) (*message.RPCReply, error) {

	// get XML payload
	request, err := marshall(operation)
	if err != nil {
		return nil, err
	}

	// setup and register callback
	reply := make(chan message.RPCReply, 1)
	callback := func(event Event) {
		reply <- *event.RPCReply()
		session.logger.Info("Successfully executed RPC")
	}
	session.Listener.Register(operation.GetMessageID(), callback)

	// send rpc
	session.logger.Info("Sending RPC")
	err = session.Transport.Send(request)
	if err != nil {
		return nil, err
	}

	select {
	case res := <-reply:
		return &res, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, errors.New("timeout while executing request")
	}
}

func marshall(operation interface{}) ([]byte, error) {
	request, err := xml.Marshal(operation)
	if err != nil {
		return nil, err
	}

	header := []byte(xml.Header)
	request = append(header, request...)
	return request, nil
}
