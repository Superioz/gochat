package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Represents this servers with its connected clients
// and the associated channels
type Server struct {
	Clients        map[uint16]ServerClient
	Port           string
	NewConnection  chan net.Conn
	DeadConnection chan net.Conn
	Messages       chan Message
	Listener       net.Listener
}

// Represents a connected client to this server
// with the id representing the index inside the
// client map - 1
type ServerClient struct {
	Id         uint16
	Name       string
	Connection net.Conn
	Passed     bool
}

// Represents a single message sent by a client
type Message struct {
	Client  ServerClient
	Message string
}

// Creates a new server instance
func NewServer(p string) Server {
	return Server{
		Clients:       make(map[uint16]ServerClient),
		Port:          p,
		NewConnection: make(chan net.Conn),
		Messages:      make(chan Message),
	}
}

// Adds a new client instance to the server
func (s *Server) Add(conn net.Conn) ServerClient {
	c := ServerClient{Id: uint16(len(s.Clients)) + 1, Connection: conn}
	s.Clients[c.Id] = c
	return c
}

// Removes a client instance from the server
// For example if the connection died.
func (s *Server) Remove(cl ServerClient) {
	delete(s.Clients, cl.Id)
}

// Starts the tcp server and listens for incoming connections
// Also it handles incoming messages and timeouts
func (s *Server) ListenAndHandle() {
	server, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatal(err)
	}
	s.Listener = server

	go serverListen(s)
	go serverHandle(s)
}

// Listens to incoming connections
func serverListen(s *Server) {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// write the new connection to channel
		s.NewConnection <- conn
	}
}

// Handles the channels associated with this server
// For example if a new connection comes in it automatically
// handles the incoming messages from this instance
func serverHandle(s *Server) {
	for {
		select {
		case c := <-s.NewConnection:
			cl := s.Add(c)
			fmt.Printf("Client connected (#%d) Not passed yet though\n", cl.Id)

			// start go routine for handling incoming messages
			go serverHandleMessages(s, cl)
			break
		case m := <-s.Messages:
			// print client message to console
			fullMessage := fmt.Sprintf("Client #%d > %s", m.Client.Id, m.Message)
			fmt.Println(fullMessage)

			// broadcast the incoming message to every client
			for _, cl := range s.Clients {
				go func(c ServerClient) {
					p := NewMessagePacket(c.Id, fullMessage)
					_, err := c.Connection.Write(p.encode())

					if err != nil {
						s.DeadConnection <- c.Connection
					}
				}(cl)
			}
			break
		}
	}
}

// Handles the message of a single client by scanning
// the connection to the tcp server and adding them
// to the incoming message channel
func serverHandleMessages(s *Server, c ServerClient) {
	reader := bufio.NewReader(c.Connection)

loop:
	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("Error while reading line", err)
			break
		}

		m, err := DecodeBytes(b)

		if err != nil {
			break
		}

		switch m.(type) {
		case *HandshakePacket:
			p := m.(*HandshakePacket)
			fmt.Println("Received handshake from", p.Client, " ..")

			if !p.Passed {
				p.Passed = true
				p.ClientId = c.Id

				fmt.Println("Sending response to", p.Client, "( Passed:", p.Passed, ") ...")

				_, err := c.Connection.Write(p.encode())
				if err != nil {
					break loop
				}
			}
			break
		case *MessagePacket:
			p := m.(*MessagePacket)

			s.Messages <- Message{c, p.Message}
			break
		}
	}

	// If there is nothing to read, something must have failed
	// therefore remove the connection
	s.Remove(c)
	fmt.Printf("Client disconnected (#%d)\n", c.Id)
}
