package network

import (
	"bytes"
)

var packetRegistry = make(map[uint16]*Packet)

// Returns an abstract packet. Only decodes and encodes
// a []byte
type Packet interface {
	encode() []byte
	decode([]byte)
}

func GetPacket(id uint16) *Packet {
	return packetRegistry[id]
}

// Represents a packet which is for authentication at the server
// Also a `Client` id as string can be passed.
// `Passed` returns the result of the authentication.
// The `Id` is the id of the packet in the registry.
type HandshakePacket struct {
	Id     uint16
	Client string
	Passed bool
}

func (p *HandshakePacket) encode() []byte {
	buf := bytes.NewBuffer([]byte{})

	WriteUint16(buf, p.Id)
	WriteString(buf, p.Client)
	WriteBool(buf, p.Passed)
	return buf.Bytes()
}

func (p *HandshakePacket) decode(b []byte) {
	buf := bytes.NewBuffer(b)

	p.Id = ReadUint16(buf)
	p.Client = ReadString(buf)
	p.Passed = ReadBool(buf)
}

// Represents a packet which is for sending
// a specific message from a specific client.
// The `Id` is the id of the packet in the registry.
type MessagePacket struct {
	Id       uint16
	ClientId uint16
	Message  string
}

func (p *MessagePacket) encode() []byte {
	buf := bytes.NewBuffer([]byte{})

	WriteUint16(buf, p.Id)
	WriteUint16(buf, p.ClientId)
	WriteString(buf, p.Message)
	return buf.Bytes()
}

func (p *MessagePacket) decode(b []byte) {
	buf := bytes.NewBuffer(b)

	p.Id = ReadUint16(buf)
	p.ClientId = ReadUint16(buf)
	p.Message = ReadString(buf)
}
