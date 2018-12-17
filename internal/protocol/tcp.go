package protocol

import (
	"bufio"
	"fmt"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/nickname"
	"log"
	"net"
	"reflect"
)

type TCPClient struct {
	Nickname         string
	Connection       *net.Conn
	outgoingMessages chan *network.MessagePacket
	incomingMessages chan *network.MessagePacket
	StateUpdates     chan bool
}

func NewTCPClient() TCPClient {
	return TCPClient{Nickname: nickname.GetRandom(), outgoingMessages: make(chan *network.MessagePacket),
		incomingMessages: make(chan *network.MessagePacket), StateUpdates: make(chan bool)}
}

func (p *TCPClient) Connect(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	p.Connection = &conn

	select {
	case p.StateUpdates <- true:
	default:
	}

	go func(p *TCPClient) {
		for {
			select {
			case s := <-p.outgoingMessages:
				_, err := (*p.Connection).Write(s.Encode())
				if err != nil {
					log.Fatal(err)
				}
				break
			case m := <-p.incomingMessages:
				fmt.Println(m.Message)
				break
			}
		}
	}(p)

	go func(p *TCPClient) {
		r := bufio.NewReader(*p.Connection)

		for {
			b, _, _ := r.ReadLine()
			m, err := network.DecodeBytes(b)

			if err != nil || (reflect.TypeOf(m) != reflect.TypeOf(&network.MessagePacket{})) {
				continue
			}

			p.incomingMessages <- m.(*network.MessagePacket)
		}
	}(p)
}

func (p *TCPClient) Disconnect() error {
	err := (*p.Connection).Close()
	if err != nil {
		return err
	}
	p.StateUpdates <- false

	return nil
}

func (p TCPClient) Send(packet network.MessagePacket) {
	select {
	case p.outgoingMessages <- &packet:
		// successful
		break
	default:
		// not successful
		break
	}
}

func (p TCPClient) Receive() chan *network.MessagePacket {
	return p.incomingMessages
}

func (p TCPClient) State() chan bool {
	return p.StateUpdates
}
