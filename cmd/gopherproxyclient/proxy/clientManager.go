package proxy

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
)

type ClientManager struct {
	Client       *proxy.ProxyClient
	StateManager *stateManager
	Initialized  bool
}

// ============================================
// Constructors
// ============================================

// NewClientManager creates a new client manager
func NewClientManager(client *proxy.ProxyClient) *ClientManager {
	return &ClientManager{
		Client:       client,
		StateManager: NewStateManager(),
		Initialized:  false,
	}
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
	}
}

// ============================================
// Event Handlers
// ============================================

func (manager *ClientManager) handleData(client *proxy.ProxyClient, packet proxcom.Packet) {
	fmt.Printf("Received data packet from %s: %s\n", packet.Source.Name, string(packet.Data))
}

func (manager *ClientManager) handleError(client *proxy.ProxyClient, packet proxcom.Packet) {
	logging.Get().Errorw("Received error packet",
		"error", string(packet.Data))

}

func (manager *ClientManager) handleCriticalError(client *proxy.ProxyClient, packet proxcom.Packet) {
	logging.Get().Errorw("Received critical error packet",
		"error", string(packet.Data))
	client.Close()
	os.Exit(1)
}

func (manager *ClientManager) handleSocketConnect(client *proxy.ProxyClient, packet proxcom.Packet) {
	logging.Get().Infow("Received socket connect packet",
		"endpoint", packet.Target)
}

func (manager *ClientManager) handleSocketDisconnect(client *proxy.ProxyClient, packet proxcom.Packet) {
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
			case proxcom.Data:
				manager.handleData(client, packet)
			case proxcom.Error:
				manager.handleError(client, packet)
			case proxcom.CriticalError:
				manager.handleCriticalError(client, packet)
			case proxcom.ChannelState:
				manager.StateManager.handleChannelState(client, packet)
			case proxcom.SocketConnect:
				manager.handleSocketConnect(client, packet)
			case proxcom.SocketDisconnect:
				manager.handleSocketDisconnect(client, packet)
			}
		}
	}
}
