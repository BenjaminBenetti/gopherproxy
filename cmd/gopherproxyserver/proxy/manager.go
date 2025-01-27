package proxy

import (
	"errors"
	"sync"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/CanadianCommander/gopherproxy/internal/proxy"
	proxylib "github.com/CanadianCommander/gopherproxy/internal/proxy"
	"github.com/google/uuid"
)

type manager struct {
	clients        map[string][]*Client
	clientsMutex   sync.Mutex
	socketChannels map[string][]*SocketChannel
	socketMutex    sync.Mutex
}

var Manager = manager{
	clients:        make(map[string][]*Client),
	socketChannels: make(map[string][]*SocketChannel),
}

// ============================================
// Public Methods
// ============================================

// AddEndpoint adds a new endpoint to the proxy manager
func (manager *manager) AddEndpoint(endpoint *proxylib.ProxyClient) error {
	manager.clientsMutex.Lock()
	defer manager.clientsMutex.Unlock()

	newClient := NewClient(endpoint, nil)

	if !manager.checkChannelPasswords(endpoint.Settings.Channel, endpoint.Settings.Password) {
		return proxylib.NewAuthenticationError("Invalid password for channel: " + endpoint.Settings.Channel)
	}

	logging.Get().Infow("Adding new endpoint to manager", "channel", endpoint.Settings.Channel, "name", endpoint.Settings.Name, "id", endpoint.Id)

	if manager.clients[endpoint.Settings.Channel] == nil {
		manager.clients[endpoint.Settings.Channel] = make([]*Client, 0)
	}
	manager.clients[endpoint.Settings.Channel] = append(manager.clients[endpoint.Settings.Channel], newClient)

	go newClient.ListenForIncomingPackets()
	go manager.watchForClientClose(endpoint)
	sendStatusUpdateToChannel(manager.clients[endpoint.Settings.Channel])
	return nil
}

// RemoveEndpoint removes an endpoint from the proxy manager
// @param channel: the channel name of the endpoint to remove
// @param id: the id of the endpoint to remove
func (manager *manager) RemoveEndpoint(channel string, id uuid.UUID) error {
	manager.clientsMutex.Lock()
	defer manager.clientsMutex.Unlock()

	if manager.clients[channel] == nil {
		return errors.New("Attempted to remove endpoint from non-existent channel: " + channel)
	}

	for i, endpoint := range manager.clients[channel] {
		if endpoint.Id == id {
			logging.Get().Infow("Removing endpoint from manager", "channel", channel, "id", id)
			manager.clients[channel] = append(manager.clients[channel][:i], manager.clients[channel][i+1:]...)
		}
	}

	sendStatusUpdateToChannel(manager.clients[channel])
	return nil
}

// EstablishNewChannel establishes a new channel between two clients
// @param client: the client that is establishing the channel
// @param chanCreatePacket: the packet containing the channel information
func (manager *manager) EstablishNewChannel(client *Client, chanCreatePacket *proxcom.CreateSocketChannelPacket) {
	manager.socketMutex.Lock()
	defer manager.socketMutex.Unlock()
	logging.Get().Infow("Establishing new channel", "client", client.Id, "packet", chanCreatePacket)

	// assign the channel id and repack
	chanCreatePacket.Id = uuid.NewString()
	newPacket, err := proxy.NewPacketFromStruct(&chanCreatePacket, proxy.SocketConnect)
	if err != nil {
		logging.Get().Errorw("Failed to repack socket connect packet", "error", err)
		return
	}

	// find the sink client
	var sinkClient *Client = nil
	for _, chanClient := range manager.clients[client.ProxyClient.Settings.Channel] {
		if chanClient.MemberInfo.Id == chanCreatePacket.Sink.Id {
			sinkClient = chanClient
			break
		}
	}

	// find the source client
	var sourceClient *Client = nil
	for _, chanClient := range manager.clients[client.ProxyClient.Settings.Channel] {
		if chanClient.MemberInfo.Id == chanCreatePacket.Source.Id {
			sourceClient = chanClient
			break
		}
	}

	// save the new channel
	if manager.socketChannels[client.ProxyClient.Settings.Channel] == nil {
		manager.socketChannels[client.ProxyClient.Settings.Channel] = make([]*SocketChannel, 0)
	}
	manager.socketChannels[client.ProxyClient.Settings.Channel] = append(manager.socketChannels[client.ProxyClient.Settings.Channel], &SocketChannel{
		Id:          chanCreatePacket.Id,
		Source:      sourceClient,
		Sink:        sinkClient,
		Initialized: false,
	})

	// send the new channel to the sink
	sinkClient.ProxyClient.Write(*newPacket)
}

// FinalizeChannel finalizes the channel creation process
// @param client: the client that is finalizing the channel
// @param channel: the channel that is being finalized
// @param chanCreatePacket: the packet that was used to create the channel
func (manager *manager) FinalizeChannel(client *Client, channel *SocketChannel, chanCreatePacket *proxcom.CreateSocketChannelPacket) {
	logging.Get().Infow("Finalizing Channel! sink reports channel creation success", "client", client.Id, "channel", channel.Id)
	channel.Initialized = true

	// send the channel creation success to the source
	sourcePacket, err := proxy.NewPacketFromStruct(&chanCreatePacket, proxy.SocketConnect)
	if err != nil {
		logging.Get().Errorw("Failed to repack socket connect packet", "error", err)
		return
	}
	channel.Source.ProxyClient.Write(*sourcePacket)
}

// ============================================
// Event Handlers
// ============================================

// handleData handles data packets received from clients
func (manager *manager) HandleData(client *Client, packet *proxylib.Packet) {
	manager.clientsMutex.Lock()
	defer manager.clientsMutex.Unlock()

	for _, channel := range manager.socketChannels[client.ProxyClient.Settings.Channel] {
		if channel.Id == packet.Chan.Id && channel.Initialized {
			if client.Id == channel.Source.MemberInfo.Id {
				channel.Sink.ProxyClient.Write(*packet)
			} else {
				channel.Source.ProxyClient.Write(*packet)
			}
			return
		}
	}

	logging.Get().Warnw("Server received data packet for unknown channel", "client", client.Id, "packet", packet)
}

// handleError handles error packets received from clients
func (manager *manager) HandleError(client *Client, packet *proxylib.Packet) {
	var err error
	packet.DecodeJsonData(&err)

	logging.Get().Errorw("Server received error from client", "client", client.Id, "error", err.Error())
}

// handleCriticalError handles critical error packets received from clients
func (manager *manager) HandleCriticalError(client *Client, packet *proxylib.Packet) {
	var err error
	packet.DecodeJsonData(&err)

	logging.Get().Errorw("Server received critical error from client", "client", client.Id, "error", err.Error())
}

// handleChannelState handles channel state packets received from clients
func (manager *manager) HandleChannelState(client *Client, packet *proxylib.Packet) {
	logging.Get().Warnw("Server received channel state packet, invalid operation", "client", client.Id)
}

func (manager *manager) HandleMemberInfo(client *Client, packet *proxylib.Packet) {
	var channelMember proxcom.ChannelMember

	err := packet.DecodeJsonData(&channelMember)
	if err != nil {
		logging.Get().Errorw("Failed to decode member info packet", "error", err)
	} else {
		logging.Get().Infow("Received new member info!", "client", client.Id)
		client.MemberInfo = &channelMember
		sendStatusUpdateToChannel(manager.clients[client.ProxyClient.Settings.Channel])
	}
}

// handleSocketConnect handles socket connect packets received from clients
func (manager *manager) HandleSocketConnect(client *Client, packet *proxylib.Packet) {
	manager.socketMutex.Lock()
	logging.Get().Infow("Received socket connect packet", "client", client.Id, "packet", packet)

	// decode the packet
	chanCreatePacket := proxcom.CreateSocketChannelPacket{}
	err := packet.DecodeJsonData(&chanCreatePacket)
	if err != nil {
		logging.Get().Errorw("Failed to decode socket connect packet", "error", err)
		manager.socketMutex.Unlock()
		return
	}

	// check if this packet is for an existing channel
	if manager.socketChannels[client.ProxyClient.Settings.Channel] != nil {
		for _, channel := range manager.socketChannels[client.ProxyClient.Settings.Channel] {
			if channel.Id == chanCreatePacket.Id && !channel.Initialized {
				manager.socketMutex.Unlock()
				manager.FinalizeChannel(client, channel, &chanCreatePacket)
				return
			}
		}
	}

	// else create a new channel
	manager.socketMutex.Unlock()
	manager.EstablishNewChannel(client, &chanCreatePacket)
}

// handleSocketDisconnect handles socket disconnect packets received from clients
func (manager *manager) HandleSocketDisconnect(client *Client, packet *proxylib.Packet) {
	manager.socketMutex.Lock()
	defer manager.socketMutex.Unlock()

	// decode the packet
	disconnectPacket := proxcom.DisconnectSocketChannelPacket{}
	err := packet.DecodeJsonData(&disconnectPacket)
	if err != nil {
		logging.Get().Errorw("Failed to decode socket disconnect packet", "error", err)
		return
	}

	// Cleanup socket channel
	if manager.socketChannels[client.ProxyClient.Settings.Channel] != nil {
		for idx, channel := range manager.socketChannels[client.ProxyClient.Settings.Channel] {
			if channel.Id == disconnectPacket.Id {
				logging.Get().Infow("Closing socket channel", "channel", channel.Id, "client", client.Id)

				// notify the other client that the channel is closing
				if client.Id == channel.Source.MemberInfo.Id {
					channel.Sink.ProxyClient.Write(*packet)
				} else {
					channel.Source.ProxyClient.Write(*packet)
				}

				// remove the channel
				manager.socketChannels[client.ProxyClient.Settings.Channel] = append(manager.socketChannels[client.ProxyClient.Settings.Channel][:idx], manager.socketChannels[client.ProxyClient.Settings.Channel][idx+1:]...)
				return
			}
		}
	}

}

// ============================================
// Go Routines
// ============================================

func (manager *manager) watchForClientClose(proxyClient *proxylib.ProxyClient) {
	<-proxyClient.CloseChannel
	Manager.RemoveEndpoint(proxyClient.Settings.Channel, proxyClient.Id)
}

// ============================================
// Private Methods
// ============================================

// checkChannelPasswords checks if the given password is valid for the given channel
// a password is valid if all other clients in that channel have the same password
// @param channel: the channel to check the password for
// @param password: the password to check
func (manager *manager) checkChannelPasswords(channel string, password string) bool {
	if manager.clients[channel] != nil {
		for _, client := range manager.clients[channel] {
			if client.ProxyClient.Settings.Password != password {
				return false
			}
		}
	}
	return true
}
