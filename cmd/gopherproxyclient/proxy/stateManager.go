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
// Public Methods
// ============================================

// SendOurChannelMemberInfoToServer sends the channel member info to the server
// for this client. This lets the server know about our current state.
func (manager *stateManager) SendOurChannelMemberInfoToServer() error {
	channelMember := manager.ClientManager.GetChannelMemberInfo()

	packet, err := proxy.NewPacketFromStruct(channelMember, proxy.MemberInfo)
	if err != nil {
		logging.Get().Errorw("Failed to create member info packet", "error", err)
		return err
	}

	manager.ClientManager.Client.Write(*packet)
	return nil
}

// ============================================
// Event Handlers
// ============================================

func (manager *stateManager) handleChannelState(client *proxy.ProxyClient, packet proxy.Packet) {
	var channelState proxcom.ChannelStateInfo

	err := packet.DecodeJsonData(&channelState)
	if err != nil {
		logging.Get().Errorw("Channel state update failed. Error decoding state packet", "error", err)
	}

	client.Id = channelState.YourId
	manager.ChannelMembers = channelState.CurrentMembers
	manager.updateForwardingRuleValidity()

	logging.Get().Infow("Channel state updated", "channel", client.Settings.Channel, "members", len(manager.ChannelMembers))

	if !manager.Initialized {
		manager.Initialized = true
		manager.InitializationChan <- true
		close(manager.InitializationChan)
	}
}

// getChannelMemberForRule returns the channel member matching the given rule
// or nil if no such member exists
// @param rule the rule to match
// @return the channel member or nil
func (stateMan *stateManager) getChannelMemberForRule(rule *proxcom.ForwardingRule) *proxcom.ChannelMember {
	for _, member := range stateMan.ChannelMembers {
		if member.Name == rule.RemoteClient {
			return member
		}
	}
	return nil
}

func (stateMan *stateManager) updateForwardingRuleValidity() {
	stateChange := false
	for _, rule := range stateMan.ClientManager.ForwardingRules {
		matched := false
		for _, member := range stateMan.ChannelMembers {
			if rule.RemoteClient == member.Name {
				matched = true
				break
			}
		}
		if rule.Valid != matched {
			stateChange = true
		}
		rule.Valid = matched
	}

	if stateChange {
		stateMan.SendOurChannelMemberInfoToServer()
	}
}
