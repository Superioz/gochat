package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"reflect"
)

// Represents the this client with the connection
// to the tcp server and message channels
type Client struct {
	Id              uint16
	Name            string
	Connection      net.Conn
	OutgoingPackets chan Packet
	IncomingPackets chan Packet
	Passed          bool
}

// Creates a new client object
func NewClient(n string) Client {
	return Client{Name: n, OutgoingPackets: make(chan Packet), IncomingPackets: make(chan Packet)}
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
		r := bufio.NewReader(cl.Connection)

		for {
			b, _, _ := r.ReadLine()
			p, err := DecodeBytes(b)

			if err != nil {
				continue
			}

			cl.IncomingPackets <- p
		}
	}()

	// Send handshake
	if !cl.Passed {
		fmt.Println("Sending handshake to server ...")
		cl.OutgoingPackets <- NewHandshakePacket(cl.Name)
	}
}

// Listens for the channels of outgoing and
// incoming messages. If an outgoing message
// is inside the channel, it will send it to the tcp server
func clientListen(cl *Client) {
	for {
		select {
		case s := <-cl.OutgoingPackets:
			if !cl.Passed && reflect.TypeOf(s) != reflect.TypeOf(&HandshakePacket{}) {
				fmt.Println("Can't send packet if client is not passed!")
				break
			}

			_, err := cl.Connection.Write(s.encode())
			if err != nil {
				log.Fatal(err)
			}
			break
		case m := <-cl.IncomingPackets:
			switch m.(type) {
			case *HandshakePacket:
				p := m.(*HandshakePacket)
				if cl.Passed || !p.Passed {
					break
				}

				cl.Passed = p.Passed
				cl.Id = p.ClientId
				fmt.Printf("Passed connection to server with id #%d\n", cl.Id)
				break
			case *MessagePacket:
				p := m.(*MessagePacket)

				fmt.Println(p.Message)
				break
			}
			break
		}
	}
}
