package main

import (
	"fmt"
	"github.com/superioz/gochat/internal/env"
	"github.com/superioz/gochat/internal/input"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/protocol"
)

func main() {
	// Initializes the default packets
	network.InitializeRegistry()

	// creates a new client
	cl := protocol.GetClient()
	fmt.Printf("Starting %s client..\n", env.GetChatType())
	go cl.Connect(env.GetServerIp("6000"))

	// TODO maybe with console commands specify the current client to use?
	/*cl := protocol.NewAMQPClient()
	go cl.Connect("amqp://guest:guest@localhost:5672")*/

	// listens to console input for message sending
	i := input.ListenToConsole()
	for {
		select {
		case m := <-i:
			// send the input of the console to the server
			cl.Send(*network.NewMessagePacket(cl.Nickname() + ": " + string(m)))
			break
		}
	}

}
