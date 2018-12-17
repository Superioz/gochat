package protocol

import (
	"bufio"
	"fmt"
	"github.com/superioz/gochat/internal/network"
	"log"
	"net"
)

type TCPServer struct {
	Listener         net.Listener
	outgoingMessages chan *network.MessagePacket
	incomingMessages chan *network.MessagePacket
	newConnections   chan net.Conn
	deadConnections  chan net.Conn
	stateUpdates     chan bool
	Clients          map[net.Conn]uint16
}

func NewTCPServer() TCPServer {
	return TCPServer{
		outgoingMessages: make(chan *network.MessagePacket),
		incomingMessages: make(chan *network.MessagePacket),
		newConnections:   make(chan net.Conn),
		deadConnections:  make(chan net.Conn),
		stateUpdates:     make(chan bool),
		Clients:          make(map[net.Conn]uint16),
	}
}

func (s *TCPServer) Start(ip string) error {
	server, err := net.Listen("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	s.Listener = server
	fmt.Println("TCP server started. Ready for connections..")

	go func(s TCPServer) {
		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				log.Fatal(err)
			}

			// write the new connection to channel
			s.newConnections <- conn
		}
	}(*s)

	go func(s TCPServer) {
		for {
			select {
			case c := <-s.newConnections:
				s.Clients[c] = uint16(len(s.Clients) + 1)
				fmt.Printf("Client #%d connected.\n", s.Clients[c])

				// start go routine for handling incoming messages
				go func(c net.Conn) {
					reader := bufio.NewReader(c)

					for {
						b, _, err := reader.ReadLine()
						if err != nil {
							break
						}

						m, err := network.DecodeBytes(b)
						if err != nil {
							break
						}

						p := m.(*network.MessagePacket)
						s.incomingMessages <- p
					}

					s.deadConnections <- c
				}(c)
				break
			case c := <-s.deadConnections:
				id := s.Clients[c]
				delete(s.Clients, c)
				fmt.Printf("Client #%d disconnected.\n", id)
				break
			case m := <-s.incomingMessages:
				// print client message to console
				fmt.Println(m.Message)

				// broadcast the incoming message to every client
				for cl := range s.Clients {
					go func(c net.Conn, m network.MessagePacket) {
						_, err := c.Write(m.Encode())

						if err != nil {
							s.deadConnections <- c
						}
					}(cl, *m)
				}
				break
			}
		}
	}(*s)

	s.stateUpdates <- true
	return nil
}

func (s TCPServer) Stop() error {
	s.Clients = make(map[net.Conn]uint16)
	err := s.Listener.Close()
	if err != nil {
		return err
	}

	// update state
	select {
	case s.stateUpdates <- false:
	}
	return nil
}

func (s TCPServer) Send() chan *network.MessagePacket {
	return s.outgoingMessages
}

func (s TCPServer) Receive() chan *network.MessagePacket {
	return s.incomingMessages
}

func (s TCPServer) State() chan bool {
	return s.stateUpdates
}
