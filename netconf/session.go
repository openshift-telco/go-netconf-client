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
	"fmt"
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
	Transport                   Transport
	SessionID                   int
	Capabilities                []string
	IsClosed                    bool
	Listener                    *Dispatcher
	IsNotificationStreamCreated bool
}

// NewSession creates a new NETCONF session using the provided transport layer.
func NewSession(t Transport) *Session {
	s := new(Session)
	s.Transport = t

	// Receive server Hello message
	serverHello, _ := s.ReceiveHello()
	s.SessionID = serverHello.SessionID
	s.Capabilities = serverHello.Capabilities

	s.Listener = &Dispatcher{}
	s.Listener.init()

	return s
}

// SendHello send the initial message through NETCONF to advertise supported capability.
func (session *Session) SendHello(hello *message.Hello) error {
	val, err := xml.Marshal(hello)
	if err != nil {
		return err
	}

	header := []byte(xml.Header)
	val = append(header, val...)
	err = session.Transport.Send(val)

	// Set Transport version after sending hello-message,
	// so the hello-message is sent using netconf:1.0 framing
	session.Transport.SetVersion("v1.0")
	for _, capability := range session.Capabilities {
		if strings.Contains(capability, message.NetconfVersion11) {
			session.Transport.SetVersion("v1.1")
			break
		}
	}

	// FIXME shouldn't be in SendHello function
	// Once the hello-message exchange is done, start listening to incoming messages
	session.listen()

	return err
}

// ReceiveHello is the first message received when connecting to a NETCONF server.
// It provides the supported capabilities of the server.
func (session *Session) ReceiveHello() (*message.Hello, error) {
	session.IsClosed = false

	hello := new(message.Hello)

	val, err := session.Transport.Receive()
	if err != nil {
		return hello, err
	}

	err = xml.Unmarshal(val, hello)
	return hello, err
}

// Close is used to close and end a session
func (session *Session) Close() error {
	session.IsClosed = true
	return session.Transport.Close()
}

// Listen starts a goroutine that listen to incoming messages and dispatch them as then are processed.
func (session *Session) listen() {
	go func() {
		for ok := true; ok; ok = !session.IsClosed {
			rawXML, err := session.Transport.Receive()
			if err != nil {
				// What should we do here?
				continue
			}
			var rawReply = string(rawXML)
			if strings.Contains(rawReply, "<rpc-reply") {

				rpcReply, err := message.NewRPCReply(rawXML)
				if err != nil {
					println(fmt.Errorf("failed to marshall message into an RPCReply. %s", err))
					continue
				}
				session.Listener.Dispatch(rpcReply.MessageID, 0, rpcReply)

			} else if strings.Contains(rawReply, "<notification") {
				notification, err := message.NewNotification(rawXML)
				if err != nil {
					println(fmt.Printf("failed to marshall message into an Notification. %s\n", err))
					continue
				}
				// In case we are using straight create-subscription, there is no way to discern who is the owner
				// of the received notification, hence we use a default handler.
				if notification.GetSubscriptionID() == "" {
					session.Listener.Dispatch(message.NetconfNotificationStreamHandler, 1, notification)
				} else {
					session.Listener.Dispatch(notification.GetSubscriptionID(), 1, notification)
				}
			} else {
				println(fmt.Errorf(fmt.Sprintf("unknown received message: \n%s", rawXML)))
			}
		}
		println("exit receiving loop")
	}()
}
