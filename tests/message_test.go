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

package tests

import (
	"encoding/xml"
	"regexp"
	"testing"

	"github.com/vitrifi/go-netconf-client/netconf/message"
)

const (
	data = "<top xmlns=\"http://example.com/schema/1.2/config\"><users/></top>"
)

var /* const */ UUIDRegex = regexp.MustCompile("\\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b")

// StripUUID uses a REGEX to remove a UUID from a provided string.
// As each NETCONF message have a unique, generated, UUID, it cannot
// be predicted when testing, hence we simply remove it.
func StripUUID(value string) string {
	return UUIDRegex.ReplaceAllString(string(value), "")
}

func panics(doesItPanic func()) (panics bool) {
	defer func() {
		if r := recover(); r != nil {
			panics = true
		}
	}()

	doesItPanic()

	return
}

func TestInvalidXML(t *testing.T) {
	invalidXML := "<<top xmlns=\"http://example.com/schema/1.2/config\"><users/></top>"
	didPanic := panics(
		func() {
			message.ValidateXML(invalidXML, message.Filter{})
		},
	)

	// expect to panic
	if didPanic != true {
		t.FailNow()
	}
}

func TestValidXML(t *testing.T) {
	invalidXML := "<top xmlns=\"http://example.com/schema/1.2/config\"><users/></top>"
	didPanic := panics(
		func() {
			message.ValidateXML(invalidXML, message.Filter{})
		},
	)

	// expect not to panic
	if didPanic == true {
		t.FailNow()
	}
}

func TestGetWithoutFilter(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><get></get></rpc>"

	rpc := message.NewGet(message.FilterTypeSubtree, "")
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestGetWithoutFilter:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestGetWithFilter(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><get><filter type=\"subtree\"><top xmlns=\"http://example.com/schema/1.2/config\"><users/></top></filter></get></rpc>"

	rpc := message.NewGet(message.FilterTypeSubtree, data)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestGetWithFilter:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestGetWithInvalidFilter(t *testing.T) {
	didPanic := panics(
		func() {
			_ = message.NewGet("dummyFilter", data)
		},
	)

	// expect to panic
	if didPanic != true {
		t.FailNow()
	}
}

func TestGetConfigWithNoFilter(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><get-config><source><running></running></source></get-config></rpc>"

	rpc := message.NewGetConfig(message.DatastoreRunning, message.FilterTypeSubtree, "")
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestGetConfigWithNoFilter:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestGetConfigWithFilter(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><get-config><source><running></running></source><filter type=\"subtree\"><top xmlns=\"http://example.com/schema/1.2/config\"><users/></top></filter></get-config></rpc>"

	rpc := message.NewGetConfig(message.DatastoreRunning, message.FilterTypeSubtree, data)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestGetConfigWithFilter:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestEditConfig(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><edit-config><target><running></running></target><default-operation>merge</default-operation><config><top xmlns=\"http://example.com/schema/1.2/config\"><users/></top></config></edit-config></rpc>"

	rpc := message.NewEditConfig(message.DatastoreRunning, message.DefaultOperationTypeMerge, data)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestEditConfig:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestEditConfigInvalidOperation(t *testing.T) {
	didPanic := panics(
		func() {
			_ = message.NewEditConfig(message.DatastoreRunning, "dummyOps", data)
		},
	)

	// expect to panic
	if didPanic != true {
		t.FailNow()
	}
}

func TestEditConfigInvalidDatastore(t *testing.T) {
	didPanic := panics(
		func() {
			_ = message.NewEditConfig("dummyDS", message.DefaultOperationTypeMerge, data)
		},
	)

	// expect to panic
	if didPanic != true {
		t.FailNow()
	}
}

func TestLock(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><lock><target><running></running></target></lock></rpc>"

	rpc := message.NewLock(message.DatastoreRunning)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestLock:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestUnlock(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><unlock><target><running></running></target></unlock></rpc>"

	rpc := message.NewUnlock(message.DatastoreRunning)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestUnlock:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestNewValidate(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><validate><source><running></running></source></validate></rpc>"

	rpc := message.NewValidate(message.DatastoreRunning)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestNewValidate:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestNewCloseSession(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><close-session></close-session></rpc>"

	rpc := message.NewCloseSession()
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestNewCloseSession:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestNewKillSession(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><kill-session><session-id>4</session-id></kill-session></rpc>"
	rpc := message.NewKillSession("4")
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestNewKillSession:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestNewCreateSubscription(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><create-subscription xmlns=\"urn:ietf:params:xml:ns:netconf:notification:1.0\"><stream>netconf-stream</stream></create-subscription></rpc>"

	rpc := message.NewCreateSubscription("", "", "netconf-stream")
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestNewCreateSubscription:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestNewEstablishSubscription(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><establish-subscription xmlns=\"urn:ietf:params:xml:ns:yang:ietf-event-notifications\" xmlns:yp=\"urn:ietf:params:xml:ns:yang:ietf-yang-push\"><stream>yp:yang-push</stream><yp:xpath-filter>/mdt-oper:mdt-oper-data/mdt-subscriptions</yp:xpath-filter><yp:period>1000</yp:period></establish-subscription></rpc>"

	rpc := message.NewEstablishSubscription("<establish-subscription xmlns=\"urn:ietf:params:xml:ns:yang:ietf-event-notifications\" xmlns:yp=\"urn:ietf:params:xml:ns:yang:ietf-yang-push\"><stream>yp:yang-push</stream><yp:xpath-filter>/mdt-oper:mdt-oper-data/mdt-subscriptions</yp:xpath-filter><yp:period>1000</yp:period></establish-subscription>")
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestNewEstablishSubscription:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestNewCommit(t *testing.T) {
	commitMsg := "some commit message"
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><commit>" + commitMsg + "</commit></rpc>"

	rpc := message.NewCommit(commitMsg)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestNewCommit:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestNewDiscardChanges(t *testing.T) {
	discardMsg := "some discard changes message"
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><discard-changes>" + discardMsg + "</discard-changes></rpc>"

	rpc := message.NewDiscardChanges(discardMsg)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestNewDiscardChanges:\nGot:%s\nWant:\n%s", got, want)
	}
}

func TestNewRPC(t *testing.T) {
	expected := "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\" message-id=\"\"><commit></commit></rpc>"
	data := "<commit></commit>"

	rpc := message.NewRPC(data)
	output, err := xml.Marshal(rpc)
	if err != nil {
		t.Errorf(err.Error())
	}

	if got, want := StripUUID(string(output)), StripUUID(expected); got != want {
		t.Errorf("TestNewRPC:\nGot:%s\nWant:\n%s", got, want)
	}
}
