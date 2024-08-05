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
func sendStatusUpdateToChannel(channelClients []*proxy.ProxyClient) {
	var channelState = proxcom.ChannelStateInfo{}
	for _, client := range channelClients {
		channelState.CurrentMembers = append(channelState.CurrentMembers, &proxcom.ChannelMember{
			Id:   client.Id,
			Name: client.Settings.Name,
		})
	}

	for _, client := range channelClients {
		packet, err := proxcom.NewPacketFromStruct(channelState, proxcom.ChannelState)
		if err != nil {
			logging.Get().Errorw("Failed to create channel state packet. Trying to continue to other clients...", "error", err)
		} else {
			client.Write(*packet)
		}
	}
}
