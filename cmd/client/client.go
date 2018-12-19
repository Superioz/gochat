package main

import (
	"fmt"
	"github.com/superioz/gochat/internal/env"
	"github.com/superioz/gochat/internal/input"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/protocol"
	"os"
)

func main() {
	// Initializes the default packets
	network.InitializeRegistry()

	// TODO remove amqp client for testing
	_ = os.Setenv("GOCHAT_TYPE", "amqp")
	_ = os.Setenv("GOCHAT_LOGGING", "true")
	_ = os.Setenv("GOCHAT_SERVER_HOST", "amqp://guest:guest@localhost")
	_ = os.Setenv("GOCHAT_SERVER_PORT", "5672")

	// creates a new client
	fmt.Printf("Starting %s client..\n", env.GetChatType())
	fmt.Println("Logging enabled:", env.IsLoggingEnabled())

	/*cl := protocol.GetClient()
	fmt.Printf("Starting %s client..\n", env.GetChatType())
	go cl.Connect(env.GetServerIp("6000"))*/

	cl := protocol.GetClient()
	go cl.Connect(env.GetServerIp("6000"))

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
