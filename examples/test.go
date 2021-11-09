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

	// Get Session
	handleReply(s.ExecRPC(message.NewGetConfig(message.DatastoreRunning, message.FilterTypeSubtree, "")))

	// Close Session
	handleReply(s.ExecRPC(message.NewCloseSession()))

}

func handleReply(reply *message.RPCReply, err error) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reply: %+v", reply.RawReply)
}