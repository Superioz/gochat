package main

import (
	"bufio"
	"github.com/streadway/amqp"
	"github.com/superioz/gochat/internal/network"
	"os"
)

func main() {
	// Initializes the default packets
	network.InitializeRegistry()

	// just to lock the repository
	_ = amqp.ChannelError

	cl := network.NewClient("hure")
	cl.ConnectAndListen("localhost:6000")

	// Read input from console and writes it to
	// outgoing message channel
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		cl.OutgoingPackets <- network.NewMessagePacket(cl.Id, s.Text())
	}
}
