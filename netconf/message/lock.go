package message

// Lock represents the NETCONF `lock` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.5
type Lock struct {
	RPC
	Target *Datastore `xml:"lock>target"`
}

// NewLock can be used to create a `lock` message.
func NewLock(datastoreType string) *Lock {
	var rpc Lock
	rpc.Target = datastore(datastoreType)
	rpc.MessageID = uuid()
	return &rpc
}
