package proxy

import (
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	proxylib "github.com/CanadianCommander/gopherproxy/internal/proxy"
	"github.com/google/uuid"
)

type Client struct {
	Id          uuid.UUID
	ProxyClient *proxylib.ProxyClient
	MemberInfo  *proxcom.ChannelMember
}

// ============================================
// Constructors
// ============================================

// NewClient creates a new client
func NewClient(proxyClient *proxylib.ProxyClient, memberInfo *proxcom.ChannelMember) *Client {
	return &Client{
		Id:          proxyClient.Id,
		ProxyClient: proxyClient,
		MemberInfo:  memberInfo,
	}
}

// ============================================
// Public Methods
// ============================================

// listenForIncomingPackets listens for incoming packets from the client
func (client *Client) ListenForIncomingPackets() {
	for {
		packet, ok := client.ProxyClient.Read()
		if !ok {
			break
		}
		RoutePacket(&packet, client)
	}
}
