// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found here
// https://github.com/Juniper/go-netconf/blob/master/LICENSE.

// The content has been modified from the original version, but the initial code
// remains from Juniper Networks, following above licence.

package netconf

import (
	"context"
	"encoding/xml"
	"io"
	"log/slog"
	"regexp"
	"strings"

	"github.com/openshift-telco/go-netconf-client/netconf/message"
)

// DefaultCapabilities sets the default capabilities of the client library.
var DefaultCapabilities = []string{
	message.NetconfVersion10,
	message.NetconfVersion11,
}

type Logger interface {
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	InfoContext(context.Context, string, ...any)
	WarnContext(context.Context, string, ...any)
	ErrorContext(context.Context, string, ...any)
}

// SessionOption allow optional configuration for the session.
type SessionOption func(*Session)

// Session represents a NETCONF sessions with a remote NETCONF server.
type Session struct {
	Transport                   Transport
	SessionID                   int
	Capabilities                []string
	IsClosed                    bool
	Listener                    *Dispatcher
	IsNotificationStreamCreated bool
	logger                      Logger
}

// NewSession creates a new NETCONF session using the provided transport layer.
func NewSession(t Transport, options ...SessionOption) (*Session, error) {
	s := new(Session)
	for _, opt := range options {
		opt(s)
	}

	if s.logger == nil {
		s.logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}

	s.Transport = t

	// Receive server Hello message
	serverHello, err := s.ReceiveHello()
	if err != nil {
		return nil, err
	}
	s.SessionID = serverHello.SessionID
	s.Capabilities = serverHello.Capabilities

	s.Listener = &Dispatcher{}
	s.Listener.init()

	return s, nil
}

// WithSessionLogger set the session logger provided in the session option.
func WithSessionLogger(logger Logger) SessionOption {
	return func(s *Session) {
		s.logger = logger
	}
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

// Listen starts a goroutine that listen to incoming messages and dispatch them as they are processed.
func (session *Session) listen() {
	go func() {
		for ok := true; ok; ok = !session.IsClosed {
			rawXML, err := session.Transport.Receive()
			if err != nil {
				// What should we do here?
				continue
			}
			var rawReply = string(rawXML)
			isRpcReply, err := regexp.MatchString(message.RpcReplyRegex, rawReply)
			if err != nil {
				session.logger.Error("failed to match RPCReply",
					"rawReply", rawReply,
					"err", err,
				)
				continue
			}

			if isRpcReply {
				rpcReply, err := message.NewRPCReply(rawXML)
				if err != nil {
					session.logger.Error("failed to marshall message into an RPCReply",
						"err", err,
					)
					continue
				}
				session.Listener.Dispatch(rpcReply.MessageID, 0, rpcReply)
				continue
			}

			isNotification, err := regexp.MatchString(message.NotificationMessageRegex, rawReply)
			if err != nil {
				session.logger.Error("failed to match notification",
					"rawReply", rawReply,
					"err", err,
				)
				continue
			}
			if isNotification {
				notification, err := message.NewNotification(rawXML)
				if err != nil {
					session.logger.Error("failed to marshall message into an Notification",
						"err", err,
					)
					continue
				}
				// In case we are using straight create-subscription, there is no way to discern who is the owner
				// of the received notification, hence we use a default handler.
				if notification.GetSubscriptionID() == "" {
					session.Listener.Dispatch(message.NetconfNotificationStreamHandler, 1, notification)
				} else {
					session.Listener.Dispatch(notification.GetSubscriptionID(), 1, notification)
				}
				continue
			}

			session.logger.Error("unknown received message",
				"rawXML", rawXML,
			)
		}
		session.logger.Info("exit receiving loop")
	}()
}
