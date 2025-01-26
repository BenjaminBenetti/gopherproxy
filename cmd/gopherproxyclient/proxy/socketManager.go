package proxy

import (
	"errors"
	"fmt"
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
	// map, channel id -> socket
	Sockets map[string][]*net.TCPConn

	socketMutex          sync.Mutex
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
		ClientManager:        clientManager,
		Listeners:            make([]*net.TCPListener, 0),
		Sockets:              make(map[string][]*net.TCPConn),
		socketChannelCreated: make(chan proxcom.CreateSocketChannelPacket, 10),
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

// SendDataToSocket sends data in the packet to a socket based on active socket channels
// @param packet the packet to send
// @return an error if one occurred
func (socketManager *SocketManager) SendDataToSocket(packet *proxy.Packet) error {
	socketManager.socketMutex.Lock()
	defer socketManager.socketMutex.Unlock()

	if socketManager.Sockets[packet.Chan.Id] == nil {
		return errors.New("could not find a socket for the channel id")
	}

	for _, socket := range socketManager.Sockets[packet.Chan.Id] {
		_, err := socket.Write(packet.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// ConnectOutbound connects to the server on the specified port in the forwarding rule
func (socketManager *SocketManager) ConnectOutbound(socketChannel proxcom.CreateSocketChannelPacket) {
	logging.Get().Debugw("Connecting to outbound server", "channelId", socketChannel.Id, "remoteHost", socketChannel.ForwardingRule.RemoteHost, "remotePort", socketChannel.ForwardingRule.RemotePort)

	// connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", socketChannel.ForwardingRule.RemoteHost, socketChannel.ForwardingRule.RemotePort))
	if err != nil {
		logging.Get().Errorw("Error connecting to outbound server", "error", err)
		return
	}

	logging.Get().Debugw("Established outgoing socket", "channelId", socketChannel.Id)
	go socketManager.packetPump(conn.(*net.TCPConn), socketChannel.Id)

	// notify proxy server
	packet, err := proxy.NewPacketFromStruct(&socketChannel, proxy.SocketConnect)
	if err != nil {
		logging.Get().Errorw("Error notifying proxy server of successful connect ", "error", err)
		conn.Close()
		return
	}

	socketManager.AddChannelSocket(socketChannel.Id, conn.(*net.TCPConn))
	socketManager.ClientManager.Client.Write(*packet)
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
	logging.Get().Debugw("Establishing socket channel", "rule", rule)

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

// AddChannelSocket adds a socket to the socket manager linked to a channel
func (socketManager *SocketManager) AddChannelSocket(channelId string, conn *net.TCPConn) {
	socketManager.socketMutex.Lock()
	defer socketManager.socketMutex.Unlock()

	if socketManager.Sockets[channelId] == nil {
		socketManager.Sockets[channelId] = make([]*net.TCPConn, 0)
	}
	socketManager.Sockets[channelId] = append(socketManager.Sockets[channelId], conn)
}

// ============================================
// Event Handlers
// ============================================

// handleSocketConnect sent by server to finish the establishment of a socket channel
func (socketManager *SocketManager) handleSocketConnect(_ *proxy.ProxyClient, packet proxy.Packet) {
	logging.Get().Debugw("Received socket connect packet from the server", "packet", packet)
	var createPacket proxcom.CreateSocketChannelPacket

	err := packet.DecodeJsonData(&createPacket)
	if err != nil {
		logging.Get().Errorw("Error decoding socket channel create packet", "error", err)
		return
	}

	if createPacket.Source.Id != socketManager.ClientManager.Client.Id {
		// we are not the source. Establish outgoing connection
		socketManager.ConnectOutbound(createPacket)
	} else {
		logging.Get().Debugw("Server reports socket channel created!", "packet", packet)
		socketManager.socketChannelCreated <- createPacket
	}
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

			socketManager.AddChannelSocket(channelId, conn)

			logging.Get().Debugw("Established socket channel to proxy server", "channelId", channelId)
			go socketManager.packetPump(conn, channelId)
		}
	}
}

// packetPump reads packets from the socket and forwards them to the server via the socket channel
// @param socket the socket to read packets from
// @param socketChannelId the id of the socket channel to forward packets to
func (socketManager *SocketManager) packetPump(socket *net.TCPConn, socketChannelId string) {
	buffer := make([]byte, PACKET_READ_SIZE)
	for {
		// read the packet
		bytesRead, err := socket.Read(buffer)
		if err != nil {
			logging.Get().Warn("Error reading from socket. Closing connection", "error", err)
			socket.Close()
			break
		}

		// proxy the packet.
		packet := proxy.NewPacketOfBytes(buffer[:bytesRead], proxy.Data)
		packet.Chan = proxy.SocketChannel{Id: socketChannelId}
		socketManager.ClientManager.Client.Write(*packet)
	}
}
