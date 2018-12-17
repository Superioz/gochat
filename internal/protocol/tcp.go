package protocol

import (
	"bufio"
	"github.com/superioz/gochat/internal/network"
	"log"
	"net"
	"reflect"
)

type TCPClient struct {
	Connection       *net.Conn
	OutgoingMessages chan *network.MessagePacket
	IncomingMessages chan *network.MessagePacket
	StateUpdates     chan bool
}

func (p *TCPClient) Connect(ip string) error {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return err
	}
	p.Connection = &conn
	p.StateUpdates <- true

	go func(p *TCPClient) {
		for {
			select {
			case s := <-p.OutgoingMessages:
				_, err := (*p.Connection).Write(s.Encode())
				if err != nil {
					log.Fatal(err)
				}
				break
			}
		}
	}(p)

	go func(p *TCPClient) {
		r := bufio.NewReader(*p.Connection)

		for {
			b, _, _ := r.ReadLine()
			m, err := network.DecodeBytes(b)

			if err != nil || reflect.TypeOf(m) != reflect.TypeOf(network.MessagePacket{}) {
				continue
			}

			p.IncomingMessages <- m.(*network.MessagePacket)
		}
	}(p)
	return nil
}

func (p *TCPClient) Disconnect() error {
	err := (*p.Connection).Close()
	if err != nil {
		return err
	}
	p.StateUpdates <- false

	return nil
}

func (p TCPClient) Send() chan *network.MessagePacket {
	return p.OutgoingMessages
}

func (p TCPClient) Receive() chan *network.MessagePacket {
	return p.IncomingMessages
}

func (p TCPClient) State() chan bool {
	return p.StateUpdates
}
