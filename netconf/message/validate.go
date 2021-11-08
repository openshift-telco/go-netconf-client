package message

// Validate represents the NETCONF `validate` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-8.6.4.1
type Validate struct {
	RPC
	Source *Datastore `xml:"validate>source"`
}

// NewValidate can be used to create a `lock` message.
func NewValidate(datastoreType string) *Validate {
	var rpc Validate
	rpc.Source = datastore(datastoreType)
	rpc.MessageID = uuid()
	return &rpc
}
