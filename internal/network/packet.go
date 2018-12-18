package network

import (
	"bytes"
	"errors"
	"strings"
)

var packetRegistry = make(map[uint16]Packet)

// Returns an abstract packet. Only decodes and encodes
// a []byte
type Packet interface {
	Encode() []byte
	Decode([]byte) error
}

func GetPacket(id uint16) Packet {
	return packetRegistry[id]
}

func InitializeRegistry() {
	packetRegistry[0] = &MessagePacket{Id: 0}
}

func DecodeBytes(b []byte) (Packet, error) {
	buf := bytes.NewBuffer(b)
	pid, err := ReadUint16(buf)

	if err != nil {
		return nil, err
	}

	packet := GetPacket(pid)
	if packet == nil {
		return nil, errors.New("can't find packet with given id")
	}

	err = packet.Decode(b)
	return packet, err
}

// Represents a packet which is for sending
// a specific message from a specific client.
// The `Id` is the id of the packet in the registry.
type MessagePacket struct {
	Id      uint16
	Message string
}

func (p *MessagePacket) UserAndMessage() (string, string) {
	spl := strings.Split(p.Message, ":")

	if len(spl) == 1 {
		return "", spl[0]
	}
	return spl[0], strings.Trim(spl[1], " ")
}

func NewMessagePacket(m string) *MessagePacket {
	return &MessagePacket{Id: 0, Message: m}
}

func (p *MessagePacket) Encode() []byte {
	buf := bytes.NewBuffer([]byte{})

	WriteUint16(buf, p.Id)
	WriteString(buf, p.Message)
	buf.WriteByte('\n')
	return buf.Bytes()
}

func (p *MessagePacket) Decode(b []byte) error {
	buf := bytes.NewBuffer(b)

	p.Id, _ = ReadUint16(buf)
	p.Message, _ = ReadString(buf)
	return nil
}
