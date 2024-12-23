package proxy

import (
	"errors"
	"sync"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	proxylib "github.com/CanadianCommander/gopherproxy/internal/proxy"
	"github.com/google/uuid"
)

type manager struct {
	clients      map[string][]*Client
	clientsMutex sync.Mutex
}

var Manager = manager{
	clients: make(map[string][]*Client),
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

// ============================================
// Event Handlers
// ============================================

// handleData handles data packets received from clients
func (manager *manager) HandleData(client *Client, packet *proxylib.Packet) {
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
}

// handleSocketDisconnect handles socket disconnect packets received from clients
func (manager *manager) HandleSocketDisconnect(client *Client, packet *proxylib.Packet) {
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
