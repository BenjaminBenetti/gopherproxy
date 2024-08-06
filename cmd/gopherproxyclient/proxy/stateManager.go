package proxy

import (
	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
)

type stateManager struct {
	ChannelMembers []*proxcom.ChannelMember
	ClientManager  *ClientManager
	// channel will receive true when the client is fully setup and ready to go
	InitializationChan chan bool
	Initialized        bool
}

// ============================================
// Constructors
// ============================================

// NewStateManager creates a new state manager
func NewStateManager(clientManager *ClientManager) *stateManager {
	return &stateManager{
		ChannelMembers:     make([]*proxcom.ChannelMember, 0),
		ClientManager:      clientManager,
		InitializationChan: make(chan bool, 1),
		Initialized:        false,
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

	client.Id = channelState.YourId
	manager.ChannelMembers = channelState.CurrentMembers
	manager.updateFowardingRuleValidity()

	logging.Get().Infow("Channel state updated", "channel", client.Settings.Channel, "members", len(manager.ChannelMembers))

	if !manager.Initialized {
		manager.Initialized = true
		manager.InitializationChan <- true
		close(manager.InitializationChan)
	}
}

func (stateMan *stateManager) updateFowardingRuleValidity() {
	for _, rule := range stateMan.ClientManager.ForwardingRules {
		matched := false
		for _, member := range stateMan.ChannelMembers {
			if rule.RemoteClient == member.Name {
				matched = true
				break
			}
		}
		rule.Valid = matched
	}
}
