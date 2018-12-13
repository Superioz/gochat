package main

import (
	"bufio"
	"github.com/superioz/gochat/internal/network"
	"os"
)

func main() {
	cl := network.NewClient()
	cl.ConnectAndListen("localhost:6000")

	// Read input from console and writes it to
	// outgoing message channel
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := s.Text()
		cl.OutgoingMessages <- m
	}
}
