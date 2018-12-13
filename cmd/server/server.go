package main

import (
	"fmt"
	"github.com/superioz/gochat/internal/network"
)

func main() {
	s := network.NewServer(":6000")
	fmt.Println("Starting tcp server @" + s.Port + " ...")

	// Starts the tcp server
	s.ListenAndHandle()
	fmt.Println("Started tcp server.")

	// Application has to be killed to exit
	for {
		select {}
	}
}
