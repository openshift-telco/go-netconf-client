package netconf

import (
	"encoding/xml"
	"github.com/adetalhouet/go-netconf/netconf/message"

)

// ExecRPC is used to execute an RPC method
func (s *Session) ExecRPC(operation interface{}) (*message.RPCReply, error) {
	request, err := xml.Marshal(operation)
	if err != nil {
		return nil, err
	}

	header := []byte(xml.Header)
	request = append(header, request...)

	err = s.Transport.Send(request)
	if err != nil {
		return nil, err
	}

	rawXML, err := s.Transport.Receive()
	if err != nil {
		return nil, err
	}

	reply, err := message.NewRPCReply(rawXML)
	if err != nil {
		return nil, err
	}

	return reply, nil
}