package protocol

import (
	"github.com/superioz/gochat/internal/env"
	"github.com/superioz/gochat/internal/network"
)

type Client interface {
	Connect(ip string)
	Disconnect() error
	Send(packet network.MessagePacket)
	Receive() chan *network.MessagePacket
	State() chan bool

	Nickname() string
}

type Server interface {
	Start(ip string) error
	Stop() error
	Send() chan *network.MessagePacket
	Receive() chan *network.MessagePacket
	State() chan bool
}

func GetClient() Client {
	t := env.GetProtocol()

	var c Client
	if t == "amqp" {
		cl := NewAMQPClient()
		c = &cl
	} else {
		cl := NewTCPClient()
		c = &cl
	}
	return c
}
