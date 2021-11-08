package message

import "encoding/xml"

// Hello is the message sent when a new NETCONF session is established.
// https://datatracker.ietf.org/doc/html/rfc6241#section-8.1
type Hello struct {
	XMLName      xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 hello"`
	Capabilities []string `xml:"capabilities>capability"`
	SessionID    int      `xml:"session-id,omitempty"`
}