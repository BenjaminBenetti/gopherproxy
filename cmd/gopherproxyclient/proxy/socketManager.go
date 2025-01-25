package proxy

import (
	"net"
	"sync"
)

type SocketManager struct {
	ClientManager *ClientManager
	Listeners     []*net.TCPListener
	Sockets       []*net.TCPConn

	listenerMutex sync.Mutex
}

// ============================================
// Constructors
// ============================================

// NewSocketManager creates a new socket manager
func NewSocketManager(clientManager *ClientManager) *SocketManager {
	return &SocketManager{
		ClientManager: clientManager,
		Listeners:     make([]*net.TCPListener, 0),
		Sockets:       make([]*net.TCPConn, 0),
	}
}

// ============================================
// Public Methods
// ============================================

// Listen starts the socket manager listening on the specified port
// @param port the port to listen on
// @param tcpType the type of tcp to listen on, can be either "tcp" or "tcp4" or "tcp6"
func (socketManager *SocketManager) Listen(port int, tcpType string) {
	socketManager.listenerMutex.Lock()
	defer socketManager.listenerMutex.Unlock()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		panic(err)
	}

	socketManager.Listeners = append(socketManager.Listeners, listener)
}

// Close closes the socket manager. Closing all listeners and sockets
func (socketManager *SocketManager) Close() {
	socketManager.listenerMutex.Lock()
	defer socketManager.listenerMutex.Unlock()

	for _, listener := range socketManager.Listeners {
		listener.Close()
	}
}
