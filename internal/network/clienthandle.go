package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Represents the this client with the connection
// to the tcp server and message channels
type Client struct {
	Connection       net.Conn
	OutgoingMessages chan string
	IncomingMessages chan string
}

// Creates a new client object
func NewClient() Client {
	return Client{OutgoingMessages: make(chan string), IncomingMessages: make(chan string)}
}

// Connects to a tcp server with given ip
// If the connection fails, the application will stop
func (cl *Client) ConnectAndListen(ip string) {
	conn, err := net.Dial("tcp", ip)

	if err != nil {
		log.Fatal(err)
	}
	cl.Connection = conn

	go clientListen(cl)

	// Read incoming messages from the server
	// and put them into the incoming channel
	go func() {
		r := bufio.NewScanner(cl.Connection)

		for r.Scan() {
			cl.IncomingMessages <- r.Text()
		}
	}()
}

// Listens for the channels of outgoing and
// incoming messages. If an outgoing message
// is inside the channel, it will send it to the tcp server
func clientListen(cl *Client) {
	for {
		select {
		case s := <-cl.OutgoingMessages:
			_, err := cl.Connection.Write([]byte(s + "\n"))
			if err != nil {
				log.Fatal(err)
			}
			break
		case m := <-cl.IncomingMessages:
			fmt.Println(m)
			break
		}
	}
}
