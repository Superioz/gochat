package protocol

import (
	"github.com/streadway/amqp"
	"github.com/superioz/gochat/internal/network"
	"github.com/superioz/gochat/internal/nickname"
	"log"
)

const queueName string = "goqueue"

// represents an amqp client
type AMQPClient struct {
	Nickname         string
	Connection       *amqp.Connection
	outgoingMessages chan *network.MessagePacket
	incomingMessages chan *network.MessagePacket
	stateUpdates     chan bool
	Channel          *amqp.Channel
	Queue            amqp.Queue
}

func NewAMQPClient() AMQPClient {
	return AMQPClient{Nickname: nickname.GetRandom(), outgoingMessages: make(chan *network.MessagePacket),
		incomingMessages: make(chan *network.MessagePacket), stateUpdates: make(chan bool)}
}

// connects the client to the amqp server
// uses the `ip` to connect to the server
func (p *AMQPClient) Connect(ip string) {
	conn, err := amqp.Dial(ip)
	if err != nil {
		log.Fatal(err)
	}
	p.Connection = conn

	// creates a channel to amqp server
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	p.Channel = ch

	// declares a new queue with name=`queueName`
	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	p.Queue = q

	// handles outgoing messages
	go func(p *AMQPClient) {
		for {
			select {
			// publishes a new message to the amqp server
			// if the channel received a new message
			case s := <-p.outgoingMessages:
				err = ch.Publish("", q.Name, false, false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(s.Message),
					})
				break
			}
		}
	}(p)

	// get the channel consumer
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// handles incoming messages
	go func(p *AMQPClient) {
		for d := range msgs {
			// send a message packet to incoming handler
			// if a new message is inside the channel
			p.incomingMessages <- &network.MessagePacket{Message: string(d.Body)}
		}
	}(p)

	select {
	case p.stateUpdates <- true:
		break
	}
}

// disconnects the client from the amqp server
func (p *AMQPClient) Disconnect() error {
	err := p.Connection.Close()
	if err != nil {
		return err
	}

	// change state
	select {
	case p.stateUpdates <- false:
	}
	return nil
}

// sends a message packet to the server
func (p AMQPClient) Send(packet network.MessagePacket) {
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
func (p AMQPClient) Receive() chan *network.MessagePacket {
	return p.incomingMessages
}

// returns the current connection state channel
func (p AMQPClient) State() chan bool {
	return p.stateUpdates
}
