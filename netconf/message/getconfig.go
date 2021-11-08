package message

// GetConfig represents the NETCONF `get-config` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.1
type GetConfig struct {
	RPC
	Source *Datastore `xml:"get-config>source"`
	Filter *Filter    `xml:"get-config>filter"`
}

// NewGetConfig can be used to create a `get-config` message.
func NewGetConfig(datastoreType string, filterType string, filterData string) *GetConfig {
	var rpc GetConfig
	if filterData != "" {
		validateXML(filterData, Filter{})
		validateFilterType(filterType)

		filter := Filter{
			Type: filterType,
			Data: filterData,
		}
		rpc.Filter = &filter
	}
	rpc.Source = datastore(datastoreType)
	rpc.MessageID = uuid()
	return &rpc
}
