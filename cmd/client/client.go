package main

import (
	"github.com/superioz/gochat/internal/input"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/protocol"
)

func main() {
	// Initializes the default packets
	network.InitializeRegistry()

	// creates a new client
	cl := protocol.NewTCPClient()
	go cl.Connect("localhost:6000")

	// listens to console input for message sending
	i := input.ListenToConsole()
	for {
		select {
		case m := <-i:
			// send the input of the console to the server
			cl.Send(*network.NewMessagePacket(cl.Nickname + ": " + string(m)))
			break
		}
	}

}
