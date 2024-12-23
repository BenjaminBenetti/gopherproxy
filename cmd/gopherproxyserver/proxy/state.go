package proxy

import (
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
)

// ============================================
// private Methods
// ============================================

// SendStatusUpdateToChannel sends a status update to all clients in a channel
func sendStatusUpdateToChannel(channelClients []*Client) {
	var channelState = proxcom.ChannelStateInfo{}

	for _, client := range channelClients {
		if client.MemberInfo != nil {
			channelState.CurrentMembers = append(channelState.CurrentMembers, client.MemberInfo)
		}
	}

	for _, client := range channelClients {
		channelState.YourId = client.ProxyClient.Id
		packet, err := proxy.NewPacketFromStruct(channelState, proxy.ChannelState)
		if err != nil {
			logging.Get().Errorw("Failed to create channel state packet. Trying to continue to other clients...", "error", err)
		} else {
			client.ProxyClient.Write(*packet)
		}
	}
}
