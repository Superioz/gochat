package protocol

import "github.com/superioz/gochat/internal/network"

type Client interface {
	Connect(ip string) error
	Disconnect() error
	Send() chan *network.MessagePacket
	Receive() chan *network.MessagePacket
	State() chan bool
}

type Server interface {
	Start(ip string) error
	Stop() error
	Send() chan *network.MessagePacket
	Receive() chan *network.MessagePacket
	State() chan bool
}
