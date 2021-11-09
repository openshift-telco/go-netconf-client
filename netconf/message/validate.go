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

// Validate represents the NETCONF `validate` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-8.6.4.1
type Validate struct {
	RPC
	Source *Datastore `xml:"validate>source"`
}

// NewValidate can be used to create a `lock` message.
func NewValidate(datastoreType string) *Validate {
	var rpc Validate
	rpc.Source = datastore(datastoreType)
	rpc.MessageID = uuid()
	return &rpc
}
