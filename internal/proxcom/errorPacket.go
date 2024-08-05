package proxcom

// ===========================================
// Constructor
// ===========================================

func NewErrorPacket(err error) *Packet {
	return newErrorPacket(err, Error)
}

func NewCriticalErrorPacket(err error) *Packet {
	return newErrorPacket(err, CriticalError)
}

// ===========================================
// Private Methods
// ===========================================

func newErrorPacket(err error, tp PacketType) *Packet {
	return &Packet{
		Type:   tp,
		Target: Endpoint{},
		Source: Endpoint{},
		Data:   []byte(err.Error()),
	}
}
