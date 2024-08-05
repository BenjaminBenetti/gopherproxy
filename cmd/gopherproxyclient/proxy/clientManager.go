package proxy

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CanadianCommander/gopherproxy/internal/websocket"
)

type clientManager struct {
	client *websocket.ProxyClient
}

// ============================================
// Constructors
// ============================================

// NewClientManager creates a new client manager
func NewClientManager(client *websocket.ProxyClient) *clientManager {
	return &clientManager{
		client: client,
	}
}

// ============================================
// Public Methods
// ============================================

// Start starts the client manager
func (manager *clientManager) Start() {
	go messageProcessingLoop(manager.client)
	createSigtermHandler(manager.client)
}

// ============================================
// Go Routines
// ============================================

func createSigtermHandler(client *websocket.ProxyClient) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.Close()
		os.Exit(0)
	}()
}
