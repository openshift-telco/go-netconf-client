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

import (
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"io"
)

const (
	// FilterTypeSubtree represent the filter for get operation
	FilterTypeSubtree string = "subtree"
	// DatastoreRunning represents the running datastore
	DatastoreRunning string = "running"
	// DatastoreCandidate represents the candidate datastore
	DatastoreCandidate string = "candidate"
	// NetconfVersion10 is the XMLNS representing NETCONF 1.0 version
	NetconfVersion10 string = "urn:ietf:params:netconf:base:1.0"
	// NetconfVersion11 is the XMLNS representing NETCONF 1.1 version
	NetconfVersion11 string = "urn:ietf:params:netconf:base:1.1"
)

// Filter represents the filter parameter of `get` message.
// Find examples here: https://datatracker.ietf.org/doc/html/rfc6241#section-6.4
type Filter struct {
	XMLName xml.Name `xml:"filter,omitempty"`
	// Type defines the filter to use. Defaults to "subtree" and can support "XPath" if the server supports it.
	Type string      `xml:"type,attr,omitempty"`
	Data interface{} `xml:",innerxml"`
}

// Datastore represents a NETCONF data store element
type Datastore struct {
	Candidate interface{} `xml:"candidate,omitempty"`
	Running   interface{} `xml:"running,omitempty"`
}

// datastore returns a Datastore object populated with appropriate datastoreType
func datastore(datastoreType string) *Datastore {
	validateDatastore(datastoreType)
	switch datastoreType {
	case DatastoreRunning:
		return &Datastore{Running: ""}
	case DatastoreCandidate:
		return &Datastore{Candidate: ""}
	}
	return nil // should never get there
}

// uuid generates a "good enough" uuid
func uuid() string {
	b := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// validateXML checks a provided string can be properly unmarshall in the specified struct
func validateXML(data string, dataStruct interface{}) {
	err := xml.Unmarshal([]byte(data), &dataStruct)
	if err != nil {
		panic(fmt.Errorf("provided XML is not valid: %s. \n%s", data, err))
	}
}

// validateDatastore checks the provided string is a supported Datastore
func validateDatastore(datastore string) {
	switch datastore {
	case DatastoreRunning:
		return
	case DatastoreCandidate:
		return
	}
	panic(
		fmt.Errorf(
			"provided datastore is not valid: %s. Expecting `%s` or `%s`", datastore, DatastoreRunning,
			DatastoreCandidate,
		),
	)
}

// validateFilterType checks the provided string is a supported FilterType
func validateFilterType(filterType string) {
	switch filterType {
	case FilterTypeSubtree:
		return
	}
	panic(
		fmt.Errorf("provided filterType is not valid: %s. Expecting `%s`", filterType, FilterTypeSubtree),
	)
}
