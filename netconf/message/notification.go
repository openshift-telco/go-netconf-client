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
	// NetconfNotificationStreamHandler identifies the callback registration for a `create-subscription`
	NetconfNotificationStreamHandler = "DEFAULT_NOTIFICATION_STREAM"
	NotificationMessageRegex         = ".*notification"
)

// Notification defines a reply to a Notification
type Notification struct {
	XMLName   xml.Name `xml:"notification"`
	XMLNS     string   `xml:"xmlns,attr"`
	EventTime string   `xml:"eventTime"`
	EventData string   `xml:"eventData,omitempty"`
	// The ietf-yang-push model cisco is using isn't following rfc8641, hence accommodating here.
	// https://github.com/YangModels/yang/blob/master/vendor/cisco/xe/1761/ietf-yang-push.yang#L367
	SubscriptionIDCisco string `xml:"push-update>subscription-id,omitempty"`
	SubscriptionID      string `xml:"push-update>id,omitempty"`
	RawReply            string `xml:"-"`
	Data                string `xml:",innerxml"`
}

// GetSubscriptionID returns the subscriptionID
func (notification *Notification) GetSubscriptionID() string {
	if notification.SubscriptionID != "" {
		return notification.SubscriptionID
	}
	if notification.SubscriptionIDCisco != "" {
		return notification.SubscriptionIDCisco
	}
	return ""
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
	Subscription CreateSubscriptionData `xml:"create-subscription"`
}

// CreateSubscriptionData is the struct to create a `create-subscription` message
type CreateSubscriptionData struct {
	XMLNS     string `xml:"xmlns,attr"`
	Stream    string `xml:"stream,omitempty"` // default is NETCONF
	Filter    string `xml:",innerxml"`
	StartTime string `xml:"startTime,omitempty"`
	StopTime  string `xml:"stopTime,omitempty"`
}

// NewCreateSubscriptionDefault can be used to create a `create-subscription` message for the NETCONF stream.
func NewCreateSubscriptionDefault() *CreateSubscription {
	var rpc CreateSubscription
	var sub = &CreateSubscriptionData{
		NetconfNotificationXmlns, "", "", "", "",
	}
	rpc.Subscription = *sub
	rpc.MessageID = uuid()
	return &rpc
}

// NewCreateSubscription can be used to create a `create-subscription` message.
func NewCreateSubscription(stopTime string, startTime string, stream string, filter string) *CreateSubscription {
	var rpc CreateSubscription
	var sub = &CreateSubscriptionData{
		NetconfNotificationXmlns, stream, filter, startTime, stopTime,
	}
	rpc.Subscription = *sub
	rpc.MessageID = uuid()
	return &rpc
}

// EstablishSubscription represents the NETCONF `establish-subscription` message.
// https://datatracker.ietf.org/doc/html/rfc8639#section-2.4.2
// FIXME very very weak implementation: there is no validation made on the schema
type EstablishSubscription struct {
	RPC
	Data string `xml:",innerxml"`
}

// NewEstablishSubscription can be used to create a `establish-subscription` message.
func NewEstablishSubscription(data string) *EstablishSubscription {
	var rpc EstablishSubscription
	rpc.Data = data
	rpc.MessageID = uuid()
	return &rpc
}
