package main

import (
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

	// Get Config
	handleReply(s.ExecRPC(message.NewGetConfig(message.DatastoreRunning, message.FilterTypeSubtree, "")))

	// Get
	//handleReply(s.ExecRPC(message.NewGet(message.FilterTypeSubtree, "")))

	// Lock
	handleReply(s.ExecRPC(message.NewLock(message.DatastoreCandidate)))

	// EditConfig - change hostname
	data := "<native xmlns=\"http://cisco.com/ns/yang/Cisco-IOS-XE-native\"><hostname>test</hostname></native>"
	handleReply(s.ExecRPC(message.NewEditConfig(message.DatastoreCandidate, message.DefaultOperationTypeMerge, data)))

	// Commit
	handleReply(s.ExecRPC(message.NewCommit()))

	// Unlock
	handleReply(s.ExecRPC(message.NewUnlock(message.DatastoreCandidate)))

	// Close Session
	handleReply(s.ExecRPC(message.NewCloseSession()))

}

func handleReply(reply *message.RPCReply, err error) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", reply.RawReply)
}