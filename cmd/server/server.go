package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/superioz/gochat/internal/console"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/protocol"
)

// THIS service is only important if we want to use tcp connections
// or logging with `elastic-search`
func main() {
	// prints an ASCII figure to the console
	f := figure.NewFigure("GoChat - " + "TCPServer", "doom", true)
	f.Print()
	fmt.Println(" ")

	// start tcp server
	s := protocol.NewTCPServer()
	err := s.Start(":6000")
	if err != nil {
		panic(err)
	}

	i := console.ListenToConsole()
	for {
		select {
		case m := <-i:
			fmt.Println(m)
			s.Send() <- &network.MessagePacket{Message: string(m)}
			break
		}
	}
}
