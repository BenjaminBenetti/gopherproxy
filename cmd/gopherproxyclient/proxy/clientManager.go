package proxy

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
)

type ClientManager struct {
	Client          *proxy.ProxyClient
	StateManager    *stateManager
	Initialized     bool
	ForwardingRules []*proxcom.ForwardingRule
}

// ============================================
// Constructors
// ============================================

// NewClientManager creates a new client manager
func NewClientManager(client *proxy.ProxyClient, forwardingRules []*proxcom.ForwardingRule) *ClientManager {
	clientManager := ClientManager{
		Client:          client,
		Initialized:     false,
		ForwardingRules: forwardingRules,
	}
	clientManager.StateManager = NewStateManager(&clientManager)
	return &clientManager
}

// ============================================
// Public Methods
// ============================================

// Start starts the client manager
func (manager *ClientManager) Start() {
	go messageProcessingLoop(manager, manager.Client)
	createSigtermHandler(manager.Client)
}

func (manager *ClientManager) WaitForInitialization() {
	if !manager.Initialized {
		<-manager.StateManager.InitializationChan

		logging.Get().Info("Client fully initialized")
		manager.Initialized = true

		// send the initial status update
		manager.StateManager.SendOurChannelMemberInfoToServer()
	}
}

// GetChannelMemberInfo returns the channel member info for THIS client
func (manager *ClientManager) GetChannelMemberInfo() *proxcom.ChannelMember {
	return &proxcom.ChannelMember{
		Id:              manager.Client.Id,
		Name:            manager.Client.Settings.Name,
		ForwardingRules: manager.ForwardingRules,
	}
}

// AllForwardingRules returns all forwarding rules for this client and all remote clients
func (manager *ClientManager) AllForwardingRules() []*proxcom.ForwardingRule {
	var rules []*proxcom.ForwardingRule = make([]*proxcom.ForwardingRule, 0)

	for _, member := range manager.StateManager.ChannelMembers {
		if member.ForwardingRules != nil {
			for _, rule := range member.ForwardingRules {
				if !slices.Contains(rules, rule) {
					rules = append(rules, rule)
				}
			}
		}
	}

	return rules
}

// AllForwardingRulesTargetingClient returns all forwarding rules that target a specific client
// @param clientName the name of the client we are looking for rules targeting
func (manager *ClientManager) AllForwardingRulesTargetingClient(clientName string) []*proxcom.ForwardingRule {
	var rules []*proxcom.ForwardingRule = make([]*proxcom.ForwardingRule, 0)

	for _, rule := range manager.AllForwardingRules() {
		if rule.RemoteClient == clientName {
			rules = append(rules, rule)
		}
	}

	return rules
}

// AllRulesTargetingUs returns all forwarding rules that target this client
func (manager *ClientManager) AllRulesTargetingUs() []*proxcom.ForwardingRule {
	return manager.AllForwardingRulesTargetingClient(manager.Client.Settings.Name)
}

// ============================================
// Event Handlers
// ============================================

func (manager *ClientManager) handleData(client *proxy.ProxyClient, packet proxy.Packet) {
	fmt.Printf("Received data packet from %s: %s\n", packet.Source.Name, string(packet.Data))
}

func (manager *ClientManager) handleError(client *proxy.ProxyClient, packet proxy.Packet) {
	logging.Get().Errorw("Received error packet",
		"error", string(packet.Data))

}

func (manager *ClientManager) handleCriticalError(client *proxy.ProxyClient, packet proxy.Packet) {
	logging.Get().Errorw("Received critical error packet",
		"error", string(packet.Data))
	client.Close()
	os.Exit(1)
}

func (manager *ClientManager) handleSocketConnect(client *proxy.ProxyClient, packet proxy.Packet) {
	logging.Get().Infow("Received socket connect packet",
		"endpoint", packet.Target)
}

func (manager *ClientManager) handleSocketDisconnect(client *proxy.ProxyClient, packet proxy.Packet) {
	logging.Get().Infow("Received socket disconnect packet",
		"endpoint", packet.Target)
}

// ============================================
// Go Routines
// ============================================

func createSigtermHandler(client *proxy.ProxyClient) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.Close()
		os.Exit(0)
	}()
}

func messageProcessingLoop(manager *ClientManager, client *proxy.ProxyClient) {
	for {
		packet, ok := client.Read()
		if !ok {
			logging.Get().Info("Proxy connection closed")
			return
		} else {
			switch packet.Type {
			case proxy.Data:
				manager.handleData(client, packet)
			case proxy.Error:
				manager.handleError(client, packet)
			case proxy.CriticalError:
				manager.handleCriticalError(client, packet)
			case proxy.ChannelState:
				manager.StateManager.handleChannelState(client, packet)
			case proxy.SocketConnect:
				manager.handleSocketConnect(client, packet)
			case proxy.SocketDisconnect:
				manager.handleSocketDisconnect(client, packet)
			}
		}
	}
}
