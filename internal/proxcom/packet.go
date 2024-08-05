package proxcom

import (
	"bytes"
	"encoding/gob"
)

type PacketType int

// This enum indicates to the receiver what the packet is for
const (
	Data PacketType = iota
	Error
	CriticalError
	ConfigureEndpoint
	SocketConnect
	SocketDisconnect
)

type Packet struct {
	Type   PacketType
	Target Endpoint
	Source Endpoint
	Data   []byte
}

// ============================================
// Constructors
// ============================================

// FromBytes creates a new packet from a byte array
func NewPacketFromBytes(data []byte) (*Packet, error) {

	var packet = Packet{
		Type:   Data,
		Target: Endpoint{},
		Source: Endpoint{},
		Data:   nil,
	}

	var err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&packet)
	return &packet, err
}

// ============================================
// Public Methods
// ============================================

// ToBytes converts the packet to a byte array
func (packet *Packet) ToBytes() ([]byte, error) {
	var buffer bytes.Buffer

	var err = gob.NewEncoder(&buffer).Encode(packet)
	return buffer.Bytes(), err
}
