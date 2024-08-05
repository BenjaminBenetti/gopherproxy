package proxy

import (
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
)

type stateManager struct {
	channelMembers []*proxcom.ChannelMember
}

// ============================================
// Constructors
// ============================================

// NewStateManager creates a new state manager
func NewStateManager() *stateManager {
	return &stateManager{
		channelMembers: make([]*proxcom.ChannelMember, 0),
	}
}

// ============================================
// Event Handlers
// ============================================

func (manager *stateManager) handleChannelState(client *proxy.ProxyClient, packet proxcom.Packet) {
	var channelState proxcom.ChannelStateInfo

	err := packet.DecodeJsonData(&channelState)
	if err != nil {
		logging.Get().Errorw("Channel state update failed. Error decoding state packet", "error", err)
	}

	manager.channelMembers = channelState.CurrentMembers
	logging.Get().Infow("Channel state updated", "channel", client.Settings.Channel, "members", len(manager.channelMembers))
}
