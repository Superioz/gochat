package protocol

import (
	"bufio"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/superioz/gochat/internal/network"
	"log"
	"net"
	"reflect"
)

// represents a tcp client
type TCPClient struct {
	UUID             uuid.UUID
	Nick             string
	outgoingMessages chan *network.MessagePacket
	incomingMessages chan *network.MessagePacket
	stateUpdates     chan bool

	Connection *net.Conn
}

func NewTCPClient(nick string) TCPClient {
	return TCPClient{UUID: uuid.NewV4(), Nick: nick, outgoingMessages: make(chan *network.MessagePacket),
		incomingMessages: make(chan *network.MessagePacket), stateUpdates: make(chan bool)}
}

// connects the client to the tcp server
// uses the `ip` to connect to the server
func (p *TCPClient) Connect(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatal(err)
	}
	p.Connection = &conn
	fmt.Println("Connected to tcp server @" + ip + ".")

	select {
	case p.stateUpdates <- true:
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

// disconnects the client from the tcp server
func (p *TCPClient) Disconnect() error {
	err := (*p.Connection).Close()
	if err != nil {
		return err
	}
	p.stateUpdates <- false

	return nil
}

// sends a message packet to the server
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

// returns receive channel
func (p TCPClient) Receive() chan *network.MessagePacket {
	return p.incomingMessages
}

// returns the current connection state channel
func (p TCPClient) State() chan bool {
	return p.stateUpdates
}

// returns the current nickname
func (p TCPClient) Nickname() string {
	return p.Nick
}
