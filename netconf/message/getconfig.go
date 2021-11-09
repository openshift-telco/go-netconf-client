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
