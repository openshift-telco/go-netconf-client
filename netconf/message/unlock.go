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

// Unlock represents the NETCONF `unlock` message.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.6
type Unlock struct {
	RPC
	Target *Datastore `xml:"unlock>target"`
}

// NewUnlock can be used to create a `unlock` message.
func NewUnlock(datastoreType string) *Unlock {
	var rpc Unlock
	rpc.Target = datastore(datastoreType)
	rpc.MessageID = uuid()
	return &rpc
}
