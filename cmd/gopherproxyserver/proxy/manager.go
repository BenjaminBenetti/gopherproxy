package proxy

import (
	"errors"
	"sync"

	"github.com/CanadianCommander/gopherproxy/internal/logging"
	proxylib "github.com/CanadianCommander/gopherproxy/internal/proxy"
	"github.com/google/uuid"
)

type manager struct {
	endpoints     map[string][]*proxylib.ProxyClient
	endpointMutex sync.Mutex
}

var Manager = manager{
	endpoints: make(map[string][]*proxylib.ProxyClient),
}

// ============================================
// Public Methods
// ============================================

// AddEndpoint adds a new endpoint to the proxy manager
func (manager *manager) AddEndpoint(endpoint *proxylib.ProxyClient) error {
	manager.endpointMutex.Lock()
	defer manager.endpointMutex.Unlock()

	if !manager.checkChannelPasswords(endpoint.Settings.Channel, endpoint.Settings.Password) {
		return proxylib.NewAuthenticationError("Invalid password for channel: " + endpoint.Settings.Channel)
	}

	logging.Get().Infow("Adding new endpoint to manager", "channel", endpoint.Settings.Channel, "name", endpoint.Settings.Name, "id", endpoint.Id)

	if manager.endpoints[endpoint.Settings.Channel] == nil {
		manager.endpoints[endpoint.Settings.Channel] = make([]*proxylib.ProxyClient, 0)
	}
	manager.endpoints[endpoint.Settings.Channel] = append(manager.endpoints[endpoint.Settings.Channel], endpoint)

	sendStatusUpdateToChannel(manager.endpoints[endpoint.Settings.Channel])
	go watchForClientClose(endpoint)
	return nil
}

// RemoveEndpoint removes an endpoint from the proxy manager
// @param channel: the channel name of the endpoint to remove
// @param id: the id of the endpoint to remove
func (manager *manager) RemoveEndpoint(channel string, id uuid.UUID) error {
	manager.endpointMutex.Lock()
	defer manager.endpointMutex.Unlock()

	if manager.endpoints[channel] == nil {
		return errors.New("Attempted to remove endpoint from non-existent channel: " + channel)
	}

	for i, endpoint := range manager.endpoints[channel] {
		if endpoint.Id == id {
			logging.Get().Infow("Removing endpoint from manager", "channel", channel, "id", id)
			manager.endpoints[channel] = append(manager.endpoints[channel][:i], manager.endpoints[channel][i+1:]...)
		}
	}

	sendStatusUpdateToChannel(manager.endpoints[channel])
	return nil
}

func (manager *manager) GetEndpointsOnChannel(channel string) []*proxylib.ProxyClient {
	manager.endpointMutex.Lock()
	defer manager.endpointMutex.Unlock()

	return manager.endpoints[channel]
}

// ============================================
// Go Routines
// ============================================

func watchForClientClose(proxyClient *proxylib.ProxyClient) {
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
	if manager.endpoints[channel] != nil {
		for _, client := range manager.endpoints[channel] {
			if client.Settings.Password != password {
				return false
			}
		}
	}
	return true
}
