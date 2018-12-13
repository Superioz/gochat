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
	Clients        map[net.Conn]*ServerClient
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
	Connection net.Conn
}

// Represents a single message sent by a client
type Message struct {
	Client  ServerClient
	Message string
}

// Creates a new server instance
func NewServer(p string) Server {
	return Server{
		Clients:       make(map[net.Conn]*ServerClient),
		Port:          p,
		NewConnection: make(chan net.Conn),
		Messages:      make(chan Message),
	}
}

// Adds a new client instance to the server
func (s *Server) Add(conn net.Conn) *ServerClient {
	c := &ServerClient{uint16(len(s.Clients)) + 1, conn}
	s.Clients[conn] = c
	return c
}

// Removes a client instance from the server
// For example if the connection died.
func (s *Server) Remove(c net.Conn) {
	delete(s.Clients, c)
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
			fmt.Printf("Client connected (#%d)\n", cl.Id)

			// start go routine for handling incoming messages
			go serverHandleMessages(s, cl)
			break
		case m := <-s.Messages:
			// print client message to console
			fullMessage := fmt.Sprintf("Client #%d > %s", m.Client.Id, m.Message)
			fmt.Println(fullMessage)

			// broadcast the incoming message to every client
			for conn := range s.Clients {
				go func() {
					_, err := conn.Write([]byte(fullMessage + "\n"))

					if err != nil {
						s.DeadConnection <- conn
					}
				}()
			}
			break
		}
	}
}

// Handles the message of a single client by scanning
// the connection to the tcp server and adding them
// to the incoming message channel
func serverHandleMessages(s *Server, c *ServerClient) {
	reader := bufio.NewScanner(c.Connection)
	for reader.Scan() {
		incoming := reader.Text()

		s.Messages <- Message{*c, incoming}
	}

	// If there is nothing to read, something must have failed
	// therefore remove the connection
	s.Remove(c.Connection)
	fmt.Printf("Client disconnected (#%d)\n", c.Id)
}
