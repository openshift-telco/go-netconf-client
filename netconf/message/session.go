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
