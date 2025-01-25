package proxy

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
)

type SocketManager struct {
	ClientManager *ClientManager
	Listeners     []*net.TCPListener
	Sockets       []*net.TCPConn

	listenerMutex        sync.Mutex
	socketChannelCreated chan proxcom.CreateSocketChannelPacket
}

const PACKET_READ_SIZE = 1024
const SOCKET_CHANNEL_CREATE_TIMEOUT = 5 * time.Second

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
func (socketManager *SocketManager) Listen(port int, tcpType string, rule *proxcom.ForwardingRule) {
	socketManager.listenerMutex.Lock()
	defer socketManager.listenerMutex.Unlock()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		panic(err)
	}

	socketManager.Listeners = append(socketManager.Listeners, listener)

	go socketManager.listenLoop(listener, rule)
}

// Close closes the socket manager. Closing all listeners and sockets
func (socketManager *SocketManager) Close() {
	socketManager.listenerMutex.Lock()
	defer socketManager.listenerMutex.Unlock()

	for _, listener := range socketManager.Listeners {
		listener.Close()
	}
}

// EstablishSocketChannel establishes a socket channel with the server
// This channel is used to proxy packets between the source and sink defined in the forwarding rule
func (socketManager *SocketManager) EstablishSocketChannel(rule *proxcom.ForwardingRule) (string, error) {
	socketManager.listenerMutex.Lock()
	defer socketManager.listenerMutex.Unlock()

	source := *socketManager.ClientManager.GetChannelMemberInfo()
	sink := socketManager.ClientManager.StateManager.getChannelMemberForRule(rule)
	if sink == nil {
		return "", errors.New("could not find a channel member for the forwarding rule")
	}

	socketCreatePacket, newChanRequestId, err := proxcom.BuildSocketChannelCreatePacket(source, *sink, *rule)
	if err != nil {
		return "", err
	}

	// send connection request to server
	socketManager.ClientManager.Client.Write(*socketCreatePacket)

	// wait for the server to respond with the connection info
	select {
	case createPacket := <-socketManager.socketChannelCreated:
		if createPacket.RequestId == newChanRequestId {
			return createPacket.Id, nil
		}
	case <-time.After(SOCKET_CHANNEL_CREATE_TIMEOUT):
		return "", errors.New("socket channel creation timed out")
	}

	return "", errors.New("socket channel creation failed")
}

// ============================================
// Event Handlers
// ============================================

// handleSocketConnect sent by server to finish the establishment of a socket channel
func (socketManager *SocketManager) handleSocketConnect(_ *proxy.ProxyClient, packet proxy.Packet) {
	var createPacket proxcom.CreateSocketChannelPacket

	err := packet.DecodeJsonData(&createPacket)
	if err != nil {
		logging.Get().Errorw("Error decoding socket channel create packet", "error", err)
		return
	}

	socketManager.socketChannelCreated <- createPacket
}

// ============================================
// Go Routines
// ============================================

// listenLoop listens on the specified listener
func (socketManager *SocketManager) listenLoop(listener *net.TCPListener, rule *proxcom.ForwardingRule) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logging.Get().Warn("Error in TCP listener. Could not accept incomming connection. Trying to continue")
		} else {

			// establish the socket channel on server
			channelId, err := socketManager.EstablishSocketChannel(rule)
			if err != nil {
				logging.Get().Errorw("Error establishing socket channel", "error", err)
				conn.Close()
				continue
			}

			logging.Get().Infow("Established socket channel", "channelId", channelId)
			go socketManager.packetPump(conn, channelId)
		}
	}
}

// packetPump reads packets from the socket and forwards them to the server via the socket channel
// @param socket the socket to read packets from
// @param socketChannelId the id of the socket channel to forward packets to
func (socketManager *SocketManager) packetPump(socket *net.TCPConn, socketChannelId string) {
	// buffer := make([]byte, PACKET_READ_SIZE)
	// for {
	// 	// read the packet
	// 	bytesRead, err := socket.Read(buffer)
	// 	if err == nil {
	// 		logging.Get().Warn("Error reading from socket. Closing connection")
	// 		socket.Close()
	// 		break
	// 	}

	// 	// proxy the packet.
	// 	packet, err := proxy.NewPacketFromBytes(buffer[:bytesRead])
	// 	if err != nil {
	// 		logging.Get().Warn("Error creating packet from bytes. Closing connection")
	// 		break
	// 	}

	// 	socketManager.ClientManager.Client.Write()
	// }
}
