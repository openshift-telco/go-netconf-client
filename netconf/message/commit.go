package message

// Commit represents the NETCONF `commit` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-8.3.4.1
type Commit struct {
	RPC
	Commit interface{} `xml:"commit"`
}

// NewCommit can be used to create a `commit` message.
func NewCommit() *Commit {
	var rpc Commit
	rpc.Commit = ""
	rpc.MessageID = uuid()
	return &rpc
}

