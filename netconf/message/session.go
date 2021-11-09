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

package message

// CloseSession represents the NETCONF `close-session` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.8
type CloseSession struct {
	RPC
	CloseSession interface{} `xml:"close-session"`
}

// NewCloseSession can be used to create a `close-session` message.
func NewCloseSession() *CloseSession {
	var rpc CloseSession
	rpc.CloseSession = ""
	rpc.MessageID = uuid()
	return &rpc
}

// KillSession represents the NETCONF `kill-session` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.8
type KillSession struct {
	RPC
	SessionID string `xml:"kill-session>session-id"`
}

// NewKillSession can be used to create a `kill-session` message.
func NewKillSession(sessionID string) *KillSession {
	var rpc KillSession
	rpc.SessionID = sessionID
	rpc.MessageID = uuid()
	return &rpc
}
