package main

import (
	"flag"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/superioz/gochat/internal/console"
	"github.com/superioz/gochat/internal/env"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/nickname"
	"github.com/superioz/gochat/internal/protocol"
)

func main() {
	// influx
	_ = client.Message{}

	// get flags for test mode and protocol forcing
	test := flag.Bool("t", false, "sets the test environment")
	prot := flag.String("p", "", "forces a protocol")
	flag.Parse()

	// if the test flag is set, prompt to open the choose menu
	// it basically overwrites the usage of the environmental variables
	// we can use this flag for later purposes as well (maybe, idk)
	if *test {
		err := console.PromptChooseProtocol(*prot)

		if err != nil {
			panic(err)
		}
	} else if len(*prot) != 0 {
		// otherwise set the protocol flag if set
		env.SetDefaults(*prot)
	}

	nick := nickname.GetRandom()

	// prints an ASCII figure to the console
	f := figure.NewFigure("GoChat - " + nick, "doom", true)
	f.Print()
	fmt.Println(" ")

	// creates a new client
	fmt.Printf("Starting %s client..\n", env.GetProtocol())
	fmt.Println("Logging enabled:", env.IsLoggingEnabled())

	cl := protocol.GetClient(nick)
	go cl.Connect(env.GetServerIp())

	// listens to console input for message sending
	i := console.ListenToConsole()
	for {
		select {
		case m := <-i:
			// send the input of the console to the server
			cl.Send(*network.NewMessagePacket(cl.Nickname() + ": " + string(m)))
			break
		}
	}
}
