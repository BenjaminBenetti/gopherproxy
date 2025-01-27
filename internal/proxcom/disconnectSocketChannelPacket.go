package proxcom

import "github.com/CanadianCommander/gopherproxy/internal/proxy"

type DisconnectSocketChannelPacket struct {
	// channel id to disconnect
	Id string
}

// ==========================================
// Constructors
// ==========================================

// NewDisconnectSocketChannelPacket creates a new packet for disconnecting a socket channel
// The client sends this to the proxy server to request a socket channel be disconnected
// @param id: the id of the channel to disconnect
// @return: the new packet or an error if one occurred
func NewDisconnectSocketChannelPacket(id string) (*proxy.Packet, error) {
	disconnectPacket := DisconnectSocketChannelPacket{
		Id: id,
	}

	packet, err := proxy.NewPacketFromStruct(disconnectPacket, proxy.SocketDisconnect)
	if err != nil {
		return nil, err
	}
	return packet, nil
}
