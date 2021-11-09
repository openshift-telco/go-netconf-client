package message

import "fmt"

const (
	// DefaultOperationTypeMerge represents the default operation to apply when doing an edit-config operation
	DefaultOperationTypeMerge   string = "merge"
	// DefaultOperationTypeReplace represents the default operation to apply when doing an edit-config operation
	DefaultOperationTypeReplace string = "replace"
	// DefaultOperationTypeNone represents the default operation to apply when doing an edit-config operation
	DefaultOperationTypeNone    string = "none"
)

// EditConfig represents the NETCONF `edit-config` operation.
// https://datatracker.ietf.org/doc/html/rfc6241#section-7.2
type EditConfig struct {
	RPC
	Target           *Datastore `xml:"edit-config>target"`
	DefaultOperation string     `xml:"edit-config>default-operation,omitempty"`
	Config           *config    `xml:"edit-config>config"`
}

type config struct {
	Config interface{} `xml:",innerxml"`
}

// NewEditConfig can be used to create a `edit-config` message.
func NewEditConfig(datastoreType string, operationType string, data string) *EditConfig {
	validateXML(data, config{})
	validDefaultOperation(operationType)

	var rpc EditConfig
	rpc.Target = datastore(datastoreType)
	rpc.DefaultOperation = operationType
	rpc.Config = &config{Config: data}
	rpc.MessageID = uuid()
	return &rpc
}

func validDefaultOperation(operation string) {
	switch operation {
	case DefaultOperationTypeMerge:
		return
	case DefaultOperationTypeReplace:
		return
	case DefaultOperationTypeNone:
		return
	}
	panic(
		fmt.Errorf(
			"provided operation is not valid: %s. Expecting either `%s`, `%s`, or `%s`", operation,
			DefaultOperationTypeMerge, DefaultOperationTypeNone, DefaultOperationTypeReplace,
		),
	)
}
