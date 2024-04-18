package message

// CopyConfig represents the NETCONF `copy-config` operation.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.3
type CopyConfig struct {
	RPC
	Target *Datastore `xml:"copy-config>target"`
	Source *Datastore `xml:"copy-config>source"`
}

// NewCopyConfig can be used to create a `copy-config` message.
func NewCopyConfig(target string, source string) *CopyConfig {
	var rpc CopyConfig
	rpc.Target = datastore(target)
	rpc.Source = datastore(source)
	rpc.MessageID = uuid()
	return &rpc
}
