package proxcom

import (
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
	"github.com/google/uuid"
)

type CreateSocketChannelPacket struct {
	Id             string
	RequestId      string
	Source         ChannelMember
	Sink           ChannelMember
	ForwardingRule ForwardingRule
}

// ==========================================
// Public Methods
// ==========================================

// BuildSocketChannelCreatePacket builds a new packet for creating a socket channel
// The client sends this to the proxy server to request a new socket channel be setup to link the source and sink
// @param source: the source channel member
// @param sink: the sink channel member
// @param rule: the forwarding rule to use for the channel
// @return: the new packet, the request id, and an error if one occurred
func BuildSocketChannelCreatePacket(source ChannelMember, sink ChannelMember, rule ForwardingRule) (*proxy.Packet, string, error) {
	createChanPack := CreateSocketChannelPacket{
		RequestId:      uuid.NewString(),
		Source:         source,
		Sink:           sink,
		ForwardingRule: rule,
	}

	packet, err := proxy.NewPacketFromStruct(createChanPack, proxy.SocketConnect)

	if err != nil {
		return nil, "", err
	}
	return packet, createChanPack.RequestId, nil
}
