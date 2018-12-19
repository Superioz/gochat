package main

import (
	"flag"
	"fmt"
	"github.com/superioz/gochat/internal/env"
	"github.com/superioz/gochat/internal/input"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/protocol"
)

func main() {
	// get flags for test mode and protocol forcing
	test := flag.Bool("t", false, "sets the test environment")
	prot := flag.String("p", "", "forces a protocol")
	flag.Parse()

	// if the test flag is set, prompt to open the choose menu
	// it basically overwrites the usage of the environmental variables
	// we can use this flag for later purposes as well (maybe, idk)
	if *test {
		err := input.PromptChooseProtocol(*prot)

		if err != nil {
			panic(err)
		}
	} else if len(*prot) != 0 {
		// otherwise set the protocol flag if set
		env.SetDefaults(*prot)
	}

	// creates a new client
	fmt.Printf("Starting %s client..\n", env.GetProtocol())
	fmt.Println("Logging enabled:", env.IsLoggingEnabled())

	cl := protocol.GetClient()
	go cl.Connect(env.GetServerIp())

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
