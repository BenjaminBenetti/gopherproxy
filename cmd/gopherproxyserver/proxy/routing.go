package proxy

import (
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
)

// =========================================
// Public Methods
// =========================================

// RoutePacket routes a incoming packet to the correct handler based on the packet type
func RoutePacket(packet *proxy.Packet, client *Client) {
	switch packet.Type {
	case proxy.Data:
		Manager.HandleData(client, packet)
	case proxy.Error:
		Manager.HandleError(client, packet)
	case proxy.CriticalError:
		Manager.HandleCriticalError(client, packet)
	case proxy.ChannelState:
		Manager.HandleChannelState(client, packet)
	case proxy.MemberInfo:
		Manager.HandleMemberInfo(client, packet)
	case proxy.SocketConnect:
		Manager.HandleSocketConnect(client, packet)
	case proxy.SocketDisconnect:
		Manager.HandleSocketDisconnect(client, packet)
	default:
		logging.Get().Errorw("Unknown packet type", "type", packet.Type)
	}
}
