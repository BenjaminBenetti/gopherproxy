package proxy

import (
	"sync"

	"github.com/CanadianCommander/gopherproxy/internal/websocket"
)

type manager struct {
	endpoints     map[string]*websocket.ProxyClient
	endpointMutex sync.Mutex
}

var Manager = manager{
	endpoints: make(map[string]*websocket.ProxyClient),
}

// ============================================
// Public Methods
// ============================================

// AddEndpoint adds a new endpoint to the proxy manager
func (manager *manager) AddEndpoint(name string, endpoint *websocket.ProxyClient) {
	manager.endpointMutex.Lock()
	defer manager.endpointMutex.Unlock()

	manager.endpoints[name] = endpoint
}

// RemoveEndpoint removes an endpoint from the proxy manager
func (manager *manager) RemoveEndpoint(name string) {
	manager.endpointMutex.Lock()
	defer manager.endpointMutex.Unlock()

	delete(manager.endpoints, name)
}
