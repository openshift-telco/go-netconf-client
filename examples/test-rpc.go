package main

import (
	"fmt"
	"github.com/adetalhouet/go-netconf/netconf"
	"github.com/adetalhouet/go-netconf/netconf/message"
	"log"
	"sync"

	"golang.org/x/crypto/ssh"
)

func main() {

	var wg sync.WaitGroup

	// Create NETCONF session that is used to received NETCONF notification
	notificationSession := createSession()
	d := "<establish-subscription xmlns=\"urn:ietf:params:xml:ns:yang:ietf-event-notifications\" xmlns:yp=\"urn:ietf:params:xml:ns:yang:ietf-yang-push\"><stream>yp:yang-push</stream><yp:xpath-filter>/bgp-ios-xe-oper:bgp-state-data/neighbors</yp:xpath-filter><yp:period>1000</yp:period></establish-subscription>"
	handleReply(notificationSession.ExecRPC(message.NewEstablishSubscription(d)))

	wg.Add(1)
	go func() {
		defer wg.Done()
		receiveNotificationAsync(notificationSession)
	}()

	execRPC()

	// Wait for notification thread to finish
	wg.Wait()
	handleReply(notificationSession.ExecRPC(message.NewCloseSession()))
	notificationSession.Close()

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
			panic(fmt.Errorf("unknown message %s", reply))
		}
		fmt.Printf("%+v", r.RawReply)
	}
}

func execRPC() {
	// Create a second NETCONF session to perform NETCONF operations
	s := createSession()

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
	err := s.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func receiveNotificationAsync(s *netconf.Session) {
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

func createSession() *netconf.Session {
	sshConfig := &ssh.ClientConfig{
		User:            "lab",
		Auth:            []ssh.AuthMethod{ssh.Password("lab")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	s, err := netconf.DialSSH("10.64.1.54:32035", sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	capabilities := netconf.DefaultCapabilities
	err = s.SendHello(&message.Hello{Capabilities: capabilities})
	if err != nil {
		log.Fatal(err)
	}
	return s
}
