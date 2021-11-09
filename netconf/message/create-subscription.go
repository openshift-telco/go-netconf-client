package message

const (
	// NetconfNotificationXmlns is the XMLNS for the YANG model supporting NETCONF notification
	NetconfNotificationXmlns = "urn:ietf:params:xml:ns:netconf:notification:1.0"
)

// CreateSubscription represents the NETCONF `create-subscription` message.
// https://datatracker.ietf.org/doc/html/rfc5277#section-2.1.1
type CreateSubscription struct {
	RPC
	Subscription Subscription `xml:"create-subscription"`
}

// Subscription is the struct to create a `create-subscription` message
type Subscription struct {
	XMLNS     string `xml:"xmlns,attr"`
	Stream    string `xml:"stream,omitempty"`
	StartTime string `xml:"startTime,omitempty"`
	StopTime  string `xml:"stopTime,omitempty"`
}

// NewCreateSubscription can be used to create a `create-subscription` message.
func NewCreateSubscription(stopTime string, startTime string, stream string) (*CreateSubscription) {
	var rpc CreateSubscription
	var sub = &Subscription{
		NetconfNotificationXmlns, stream, startTime, stopTime,
	}
	rpc.Subscription = *sub
	rpc.MessageID = uuid()
	return &rpc
}
