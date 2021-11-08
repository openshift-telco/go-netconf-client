package netconf

import (
	"fmt"
	"github.com/adetalhouet/go-netconf/netconf/message"
	"strings"
)

// DefaultCapabilities sets the default capabilities of the client library
var DefaultCapabilities = []string{
	message.NetconfVersion10,
	message.NetconfVersion11,
}

// Session represents a NETCONF sessions with a remote NETCONF server
type Session struct {
	Transport    Transport
	SessionID    int
	Capabilities []string
}

// NewSession creates a new NETCONF session using the provided transport layer.
func NewSession(t Transport) *Session {
	s := new(Session)
	s.Transport = t

	// Receive server Hello message
	serverHello, _ := t.ReceiveHello()
	s.SessionID = serverHello.SessionID
	s.Capabilities = serverHello.Capabilities

	fmt.Println("Received server hello-message")
	fmt.Println("Sending our hello-message")

	// Send our hello using default capabilities.
	err := t.SendHello(&message.Hello{Capabilities: DefaultCapabilities})
	if err != nil {
		return nil
	}

	// Set Transport version
	t.SetVersion("v1.0")
	for _, capability := range s.Capabilities {
		if strings.Contains(capability, message.NetconfVersion11) {
			t.SetVersion("v1.1")
			break
		}
	}

	return s
}

// Close is used to close and end a session
func (s *Session) Close() error {
	return s.Transport.Close()
}

//
//// CreateSubscription is used to execute an RPC <create-subscription> operation
//func (s *Session) CreateSubscription(createSubscriptionRPC *CreateSubscriptionRPC) (*RPCReply, error) {
//	if !stringInSlice(s.ServerCapabilities, NOTIFICATION_CAPABILITY) {
//		return nil, &ErrCapabilityNotSupported{Cap: NOTIFICATION_CAPABILITY}
//	}
//	return s.send(createSubscriptionRPC, createSubscriptionRPC.MessageID)
//}
//
//
//func stringInSlice(l []string, t string) bool {
//	for _, s := range l {
//		if t == s {
//			return true
//		}
//	}
//	return false
//}
//
//type ErrCapabilityNotSupported struct {
//	Cap string
//}
//
//func (e *ErrCapabilityNotSupported) Error() string {
//	return fmt.Sprintf("capability %s is not supported", e.Cap)
//}
