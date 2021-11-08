package message

// Unlock represents the NETCONF `unlock` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.6
type Unlock struct {
	RPC
	Target *Datastore `xml:"unlock>target"`
}

// NewUnlock can be used to create a `unlock` message.
func NewUnlock(datastoreType string) *Unlock {
	var rpc Unlock
	rpc.Target = datastore(datastoreType)
	rpc.MessageID = uuid()
	return &rpc
}
