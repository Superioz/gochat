package network

import (
	"bytes"
	"errors"
)

var packetRegistry = make(map[uint16]Packet)

// Returns an abstract packet. Only decodes and encodes
// a []byte
type Packet interface {
	encode() []byte
	decode([]byte) error
}

func GetPacket(id uint16) Packet {
	return packetRegistry[id]
}

func InitializeRegistry() {
	packetRegistry[0] = &HandshakePacket{Id: 0}
	packetRegistry[1] = &MessagePacket{Id: 1}
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

	err = packet.decode(b)
	return packet, err
}

// Represents a packet which is for authentication at the server
// Also a `Client` id as string can be passed.
// `Passed` returns the result of the authentication.
// The `Id` is the id of the packet in the registry.
type HandshakePacket struct {
	Id       uint16
	Client   string
	Passed   bool
	ClientId uint16
}

func NewHandshakePacket(id string) Packet {
	return &HandshakePacket{Id: 0, Client: id, Passed: false}
}

func (p *HandshakePacket) encode() []byte {
	buf := bytes.NewBuffer([]byte{})

	WriteUint16(buf, p.Id)
	WriteString(buf, p.Client)
	WriteBool(buf, p.Passed)
	WriteUint16(buf, p.ClientId)
	buf.WriteByte('\n')
	return buf.Bytes()
}

func (p *HandshakePacket) decode(b []byte) error {
	buf := bytes.NewBuffer(b)

	p.Id, _ = ReadUint16(buf)
	p.Client, _ = ReadString(buf)
	p.Passed, _ = ReadBool(buf)
	p.ClientId, _ = ReadUint16(buf)
	return nil
}

// Represents a packet which is for sending
// a specific message from a specific client.
// The `Id` is the id of the packet in the registry.
type MessagePacket struct {
	Id       uint16
	ClientId uint16
	Message  string
}

func NewMessagePacket(id uint16, m string) Packet {
	return &MessagePacket{Id: 1, ClientId: id, Message: m}
}

func (p *MessagePacket) encode() []byte {
	buf := bytes.NewBuffer([]byte{})

	WriteUint16(buf, p.Id)
	WriteUint16(buf, p.ClientId)
	WriteString(buf, p.Message)
	buf.WriteByte('\n')
	return buf.Bytes()
}

func (p *MessagePacket) decode(b []byte) error {
	buf := bytes.NewBuffer(b)

	p.Id, _ = ReadUint16(buf)
	p.ClientId, _ = ReadUint16(buf)
	p.Message, _ = ReadString(buf)
	return nil
}
