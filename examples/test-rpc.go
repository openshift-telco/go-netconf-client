package main

import (
	"fmt"
	"github.com/adetalhouet/go-netconf/netconf"
	"github.com/adetalhouet/go-netconf/netconf/message"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

func main() {

	// Create NETCONF session
	session := createSession()

	// Define message to send
	d := "<establish-subscription xmlns=\"urn:ietf:params:xml:ns:yang:ietf-event-notifications\" xmlns:yp=\"urn:ietf:params:xml:ns:yang:ietf-yang-push\"><stream>yp:yang-push</stream><yp:xpath-filter>/bgp-ios-xe-oper:bgp-state-data/neighbors</yp:xpath-filter><yp:period>1000</yp:period></establish-subscription>"
	m := message.NewEstablishSubscription(d)

	// Define callback function for the rpc-reply
	callback := func(event netconf.Event) {
		reply := event.RPCReply()
		if reply == nil {
			println("Failed to execute RPC")
		}
		if event.EventID() == m.MessageID {
			println(fmt.Sprintf("Successfully executed Notification stream registration with subscritpionID: %s", reply.SubscriptionId))
						// if all went well, we register a callback for notification
			session.Listener.Register(
				reply.SubscriptionId, session.DefaultLogNotificationCallback(reply.SubscriptionId),
			)
		}
		println("WHAT")
	}

	// Send request
	err := session.SendRPC(m.MessageID, m, callback)
	if err != nil {
		panic(err)
	}

	execRPC(session)

	//defer session.Close()
	time.Sleep(15 * time.Second)
}

func execRPC(session *netconf.Session) {

	// Get Config
	g := message.NewGetConfig(message.DatastoreRunning, message.FilterTypeSubtree, "")
	session.SendRPC(g.MessageID, g, session.DefaultLogRpcReplyCallback(g.MessageID))

	// Get - some issues
	//handleReply(s.ExecRPC(message.NewGet(message.FilterTypeSubtree, "")))

	// Lock
	l := message.NewLock(message.DatastoreCandidate)
	session.SendRPC(l.MessageID, l, session.DefaultLogRpcReplyCallback(l.MessageID))

	// EditConfig - change hostname
	data := "<native xmlns=\"http://cisco.com/ns/yang/Cisco-IOS-XE-native\"><hostname>r1</hostname></native>"
	e := message.NewEditConfig(message.DatastoreCandidate, message.DefaultOperationTypeMerge, data)
	session.SendRPC(e.MessageID, e, session.DefaultLogRpcReplyCallback(e.MessageID))

	// Commit
	c := message.NewCommit()
	session.SendRPC(c.MessageID, c, session.DefaultLogRpcReplyCallback(c.MessageID))

	// Unlock
	u := message.NewUnlock(message.DatastoreCandidate)
	session.SendRPC(u.MessageID, u, session.DefaultLogRpcReplyCallback(u.MessageID))

	// Close Session
	d := message.NewCloseSession()
	session.SendRPC(d.MessageID, d, session.DefaultLogRpcReplyCallback(d.MessageID))
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
