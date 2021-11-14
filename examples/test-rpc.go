package main

import (
	"errors"
	"fmt"
	"github.com/adetalhouet/go-netconf/netconf"
	"github.com/adetalhouet/go-netconf/netconf/message"
	"log"

	"golang.org/x/crypto/ssh"
)

func main() {
	sshConfig := &ssh.ClientConfig{
		User:            "lab",
		Auth:            []ssh.AuthMethod{ssh.Password("lab")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	s, err := netconf.DialSSH("10.64.1.54:32035", sshConfig)

	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	fmt.Println(s.Capabilities)
	fmt.Println(s.SessionID)

	capabilities := netconf.DefaultCapabilities
	capabilities = append(capabilities, "urn:ietf:params:xml:ns:yang:ietf-event-notifications")
	s.SendHello(&message.Hello{Capabilities: capabilities})

	// Get Config
	handleReply(s.ExecRPC(message.NewGetConfig(message.DatastoreRunning, message.FilterTypeSubtree, "")))

	// Get - some issues
	//handleReply(s.ExecRPC(message.NewGet(message.FilterTypeSubtree, "")))

	// Lock
	handleReply(s.ExecRPC(message.NewLock(message.DatastoreCandidate)))

	// EditConfig - change hostname
	data := "<native xmlns=\"http://cisco.com/ns/yang/Cisco-IOS-XE-native\"><hostname>r1</hostname></native>"
	handleReply(s.ExecRPC(message.NewEditConfig(message.DatastoreCandidate, message.DefaultOperationTypeMerge, data)))

	// Commit
	handleReply(s.ExecRPC(message.NewCommit()))

	// Unlock
	handleReply(s.ExecRPC(message.NewUnlock(message.DatastoreCandidate)))

	// Close Session
	handleReply(s.ExecRPC(message.NewCloseSession()))

}

func handleReply(reply interface{}, err error) {

	if err != nil {
		panic(err)
	}

	r, ok := reply.(*message.RPCReply)
	if ok {
		fmt.Printf("%+v", r.RawReply)
	} else {
		r, ok := reply.(*message.Notification)
		if !ok {
			panic(errors.New(fmt.Sprintf("unknown message %s", reply)))
		}
		fmt.Printf("%+v", r.RawReply)
	}
}
