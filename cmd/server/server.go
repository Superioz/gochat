package main

import (
	"fmt"
	"github.com/superioz/gochat/internal/input"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/protocol"
)

// THIS service is only important if we want to use tcp connections
// or logging with `elastic-search`
func main() {
	// Initializes the default packets
	network.InitializeRegistry()

	s := protocol.NewTCPServer()
	err := s.Start(":6000")
	if err != nil {
		panic(err)
	}

	i := input.ListenToConsole()
	for {
		select {
		case m := <-i:
			fmt.Println(m)
			s.Send() <- &network.MessagePacket{Message: string(m)}
			break
		}
	}
}
