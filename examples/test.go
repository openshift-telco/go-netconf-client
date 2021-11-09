package main

import (
	"fmt"
	"github.com/adetalhouet/go-netconf/netconf"
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

	// Sends raw XML
	//reply, err := s.ExecRPC(message.NewGet("", ""))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Reply: %+v", reply)
}
