/*
Copyright 2021. Alexis de Talhouët

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package message

import "encoding/xml"

const (
	// NetconfNotificationXmlns is the XMLNS for the YANG model supporting NETCONF notification
	NetconfNotificationXmlns = "urn:ietf:params:xml:ns:netconf:notification:1.0"
)

// Notification defines a reply to a Notification
type Notification struct {
	XMLName    xml.Name `xml:"notification"`
	XMLNS      string   `xml:"xmlns,attr"`
	EventTime  string   `xml:"eventTime"`
	EventData  string   `xml:"eventData,omitempty"`
	PushUpdate string   `xml:"push-update,omitempty"`
	RawReply   string   `xml:"-"`
}

// NewNotification creates an instance of an Notification based on what was received
func NewNotification(rawXML []byte) (*Notification, error) {
	reply := &Notification{}
	reply.RawReply = string(rawXML)

	if err := xml.Unmarshal(rawXML, reply); err != nil {
		return nil, err
	}

	return reply, nil
}

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
func NewCreateSubscription(stopTime string, startTime string, stream string) *CreateSubscription {
	var rpc CreateSubscription
	var sub = &Subscription{
		NetconfNotificationXmlns, stream, startTime, stopTime,
	}
	rpc.Subscription = *sub
	rpc.MessageID = uuid()
	return &rpc
}
