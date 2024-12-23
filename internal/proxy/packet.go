package proxy

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type PacketType int

// This enum indicates to the receiver what the packet is for
const (
	Data PacketType = iota
	Error
	CriticalError
	// update all clients with the current channel state info. e.g. list of channel members
	// and other channel specific information
	ChannelState
	// sent by a channel member to update their info
	// The server will then emit a ChannelState packet to all clients in the channel
	// to let them know about the new member info
	MemberInfo
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

// NewPacketFromStruct creates a new packet from a struct with the given packet type
func NewPacketFromStruct(obj any, typ PacketType) (*Packet, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return &Packet{
		Type:   typ,
		Target: Endpoint{},
		Source: Endpoint{},
		Data:   bytes,
	}, err
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

// DecodeJsonData inside this packet.
func (packet *Packet) DecodeJsonData(out any) error {
	return json.Unmarshal(packet.Data, out)
}
