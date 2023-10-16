package tests

import (
	"encoding/xml"
	"os"
	"regexp"
	"testing"

	"github.com/openshift-telco/go-netconf-client/netconf/message"
)

func TestRPCReply(t *testing.T) {

	input, err := os.ReadFile("resources/junos-rpc-reply.xml")
	if err != nil {
		t.Fatalf("failed to read resources: %v", err)
	}

	// validate we can create RPCReply when it's encapsulated in a xmlns
	reply, err := message.NewRPCReply(input)
	if err != nil {
		t.Fatalf("failed to unmarshal rpc reply: %v", err)
	}

	// validate we can marshall the created RPCReply
	output, err := xml.Marshal(reply)
	if err != nil {
		t.Fatalf("failed to marshal rpc reply: %v", err)
	}

	if string(input) != string(output) {
		t.Errorf("got %q, \nwanted %q", string(input), string(output))
	}

	_, e := regexp.MatchString(message.RpcReplyRegex, string(input))
	if e != nil {
		t.Errorf("failed to parse rpc-reply with regex")
	}
}
