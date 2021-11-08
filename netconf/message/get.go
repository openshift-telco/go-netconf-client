package message

// Get represents the NETCONF `get` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.7
type Get struct {
	RPC
	Get          interface{} `xml:"get"`
	Filter       *Filter
}

// NewGet can be used to create a `get` message.
func NewGet(filterType string, data string) *Get {
	var rpc Get
	if data != "" {
		validateXML(data, Filter{})
		validateFilterType(filterType)

		filter := Filter{
			Type: filterType,
			Data: data,
		}
		rpc.Filter = &filter
	} else {
		rpc.Get = ""
	}
	rpc.MessageID = uuid()
	return &rpc
}
