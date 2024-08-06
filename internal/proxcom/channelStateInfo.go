package proxcom

import "github.com/google/uuid"

type ChannelStateInfo struct {
	YourId         uuid.UUID
	CurrentMembers []*ChannelMember
}
