package main

import (
	"errors"
	"fmt"
	"github.com/adetalhouet/go-netconf/netconf"
	"github.com/adetalhouet/go-netconf/netconf/message"
	"golang.org/x/crypto/ssh"
	"log"
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

	d := "<establish-subscription\n        xmlns=\"urn:ietf:params:xml:ns:yang:ietf-event-notifications\"\n        xmlns:yp=\"urn:ietf:params:xml:ns:yang:ietf-yang-push\">\n      <stream>yp:yang-push</stream>\n      <yp:xpath-filter>/mdt-oper:mdt-oper-data/mdt-subscriptions</yp:xpath-filter>\n      <yp:period>1000</yp:period>\n    </establish-subscription>"
	handleReplyN(s.ExecRPC(message.NewRPC(d)))

	condition := false
	counter := 5
	for ok := true; ok; ok = !condition {
		rawXML, err := s.Transport.Receive()
		if err != nil {
			panic(err)
		}

		var rawReply = string(rawXML)
		fmt.Printf("%+v", rawReply)

		counter--

		if counter == 0 {
			condition = true
		}
	}
}

func handleReplyN(reply interface{}, err error) {

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

