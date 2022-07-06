package udp

type LongHeader struct {
	IsLongHeader bool
	FixedBit     bool

	LongPacketType     packetType
	ReservedBits       uint8
	PacketNumberLength uint8

	Version                       VersionNumber
	DestinationConnectionIdLength uint8
	DestinationConnectionId       []byte

	SourceConnectionIdLength uint8
	SourceConnectionId       []byte

	TokenVariableLength uint32
	Token               []byte
	LengthVariable      uint32
	PacketNumber        PacketNumber
	PacketPayload       []byte

	RetryToken        []byte
	RetryIntegrityTag []byte
}

func (h *LongHeader) GetPacketNumberLength() uint8 {
	return h.PacketNumberLength + 1
}
func (h *LongHeader) GetFirstByte() byte {
	var b byte = 0
	if h.IsLongHeader {
		b = b | 0x80
	}
	if h.FixedBit {
		b = b | 0x40
	}
	b = b | (uint8(h.LongPacketType) << 4)
	b = b | (h.ReservedBits << 2)
	b = b | (h.PacketNumberLength)
	return b
}
func (h *LongHeader) SetFirstByte(b byte) {
	h.IsLongHeader = b&0x80 > 0
	h.FixedBit = b&0x40 > 0
	h.LongPacketType = packetType(b&0x30>>4)
	h.ReservedBits = 0x0c>>2
	h.PacketNumberLength = b&0x03
}

func NewLongHeader(longPacketType packetType, PlayLoad []byte, sendConfig *SendConfig) *LongHeader {
	var longHeader LongHeader
	longHeader.LongPacketType = longPacketType
	longHeader.IsLongHeader = true
	longHeader.Version = sendConfig.Version
	longHeader.DestinationConnectionIdLength = uint8(len(sendConfig.ConnectionId))
	longHeader.DestinationConnectionId = sendConfig.ConnectionId
	longHeader.SourceConnectionIdLength = 0
	longHeader.SourceConnectionId = []byte{}
	longHeader.TokenVariableLength = uint32(len(sendConfig.Token))
	longHeader.Token = sendConfig.Token
	longHeader.LengthVariable = uint32(sendConfig.PacketNumber.GetPacketNumberLength()) + uint32(len(PlayLoad))
	longHeader.PacketNumber = sendConfig.PacketNumber
	longHeader.PacketPayload = PlayLoad
	return &longHeader
}