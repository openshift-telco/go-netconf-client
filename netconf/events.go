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

package netconf

import (
	"sync"
	"time"

	"github.com/openshift-telco/go-netconf-client/netconf/message"
)

/**
This file is meant to provide all the necessary tooling to support callback mechanism.
It is to be used to subscribe listeners when NETCONF RPCs or Notifications are sent, in
order to process their response.
*/

// Names of event types
var eventTypeStrings = [...]string{
	"rpc-reply", "notification",
}

// EventType is an enumeration of the kind of events that can occur.
type EventType uint16

// String returns the name of event types
func (t EventType) String() string {
	return eventTypeStrings[t]
}

// Callback is a function that can receive events.
type Callback func(Event)

// Dispatcher objects can register callbacks for specific events, then when
// those events occur, dispatch them its according callback functions.
type Dispatcher struct {
	callbacks     map[string]Callback
	callbacksLock sync.Mutex
}

// init a dispatcher creating the callbacks map.
func (d *Dispatcher) init() {
	d.callbacks = make(map[string]Callback)
}

// Register a callback function for the specified eventID.
func (d *Dispatcher) Register(eventID string, callback Callback) {
	d.callbacksLock.Lock()
	defer d.callbacksLock.Unlock()
	d.callbacks[eventID] = callback
}

// Remove a callback function for the specified eventID.
func (d *Dispatcher) Remove(eventID string) {
	d.callbacksLock.Lock()
	defer d.callbacksLock.Unlock()
	delete(d.callbacks, eventID)
}

// WaitForMessages waits for all messages in the queue to be processed
// TODO support timeout
func (d *Dispatcher) WaitForMessages() {
	for len(d.callbacks) != 0 {
		time.Sleep(1 * time.Second)
	}
}

// Dispatch an event by triggering its associated callback.
// FIXME manage errors
func (d *Dispatcher) Dispatch(eventID string, eventType EventType, value interface{}) {
	// Create the event
	e := &event{
		eventID: eventID,
		value:   value,
	}

	// Dispatch the event to the callback
	callback := d.callbacks[eventID]
	if callback == nil {
		return
	}
	callback(e)

	// In case of rpc-reply, auto-remove registration
	// If it is a notification, we need to keep the registration active
	// as we can have still receive notification related to the subscriptionID
	switch eventType.String() {
	case "rpc-reply":
		d.Remove(eventID)
	case "notification":
		// NOOP
	}
}

// Event represents actions that occur during NETCONF exchange. Listeners can
// register callbacks with event handlers when creating a new RPC.
type Event interface {
	EventID() string
	Value() interface{}
	RPCReply() *message.RPCReply
	Notification() *message.Notification
}

// event is an internal implementation of the Event interface.
type event struct {
	eventID string
	value   interface{}
}

// EventID returns the eventID
func (e *event) EventID() string {
	return e.eventID
}

// Value returns the current value associated with the event.
func (e *event) Value() interface{} {
	return e.value
}

// RPCReply returns an RPCReply from the associated value.
func (e *event) RPCReply() *message.RPCReply {
	r, ok := e.value.(*message.RPCReply)
	if ok {
		return r
	}
	return nil
}

// Notification returns a Notification from the associated value.
func (e *event) Notification() *message.Notification {
	n, ok := e.value.(*message.Notification)
	if ok {
		return n
	}
	return nil
}
