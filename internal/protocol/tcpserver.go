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
	OutgoingMessages chan *network.MessagePacket
	IncomingMessages chan *network.MessagePacket
	NewConnections   chan net.Conn
	DeadConnections  chan net.Conn
	StateUpdates     chan bool
	Clients          map[net.Conn]uint16
}

func (s *TCPServer) Start(ip string) error {
	server, err := net.Listen("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	s.Listener = server
	s.StateUpdates <- true

	go func(s TCPServer) {
		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				log.Fatal(err)
			}

			// write the new connection to channel
			s.NewConnections <- conn
		}
	}(*s)

	go func(s TCPServer) {
		for {
			select {
			case c := <-s.NewConnections:
				s.Clients[c] = uint16(len(s.Clients)+1)

				// start go routine for handling incoming messages
				go func(c net.Conn) {
					reader := bufio.NewReader(c)

					for {
						b, _, err := reader.ReadLine()
						if err != nil {
							fmt.Println("Error while reading line", err)
							break
						}

						m, err := network.DecodeBytes(b)

						if err != nil {
							break
						}

						p := m.(*network.MessagePacket)
						s.IncomingMessages <- p
					}


				}(c)
				break
			case m := <-s.IncomingMessages:
				// print client message to console
				fmt.Println(m.Message)

				// broadcast the incoming message to every client
				for cl := range s.Clients {
					go func(c net.Conn, m network.MessagePacket) {
						_, err := c.Write(m.Encode())

						if err != nil {
							s.DeadConnections <- c
						}
					}(cl, *m)
				}
				break
			}
		}
	}(*s)

	return nil
}

func (s TCPServer) Stop() error {
	s.Clients = make(map[net.Conn]uint16)
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	s.StateUpdates <- false
	return nil
}

func (s TCPServer) Send() chan *network.MessagePacket {
	return s.OutgoingMessages
}

func (s TCPServer) Receive() chan *network.MessagePacket {
	return s.IncomingMessages
}

func (s TCPServer) State() chan bool {
	return s.StateUpdates
}
