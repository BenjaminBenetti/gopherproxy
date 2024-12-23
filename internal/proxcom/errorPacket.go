package proxcom

import "github.com/CanadianCommander/gopherproxy/internal/proxy"

// ===========================================
// Constructor
// ===========================================

func NewErrorPacket(err error) *proxy.Packet {
	return newErrorPacket(err, proxy.Error)
}

func NewCriticalErrorPacket(err error) *proxy.Packet {
	return newErrorPacket(err, proxy.CriticalError)
}

// ===========================================
// Private Methods
// ===========================================

func newErrorPacket(err error, tp proxy.PacketType) *proxy.Packet {
	return &proxy.Packet{
		Type:   tp,
		Target: proxy.Endpoint{},
		Source: proxy.Endpoint{},
		Data:   []byte(err.Error()),
	}
}
