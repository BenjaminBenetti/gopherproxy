package proxcom

import "github.com/google/uuid"

// Member of a channel. Used to keep track of who is in the channel
// and where messages can be sent.
type ChannelMember struct {
	Id   uuid.UUID
	Name string
}

// ===========================================
// Constructors
// ===========================================

// NewChannelMember creates a new channel member
// @param name: the name of the channel member
func NewChannelMember(name string) *ChannelMember {
	return &ChannelMember{
		Id:   uuid.New(),
		Name: name,
	}
}
