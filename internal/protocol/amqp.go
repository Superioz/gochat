package protocol

import (
	"github.com/streadway/amqp"
	"github.com/superioz/gochat/internal/network"
	"log"
)

const queueName string = "goqueue"

type AMQPClient struct {
	Connection       *amqp.Connection
	outgoingMessages chan *network.MessagePacket
	incomingMessages chan *network.MessagePacket
	stateUpdates     chan bool
	Channel          *amqp.Channel
	Queue            amqp.Queue
}

func (p *AMQPClient) Connect(ip string) {
	conn, err := amqp.Dial(ip)
	if err != nil {
		log.Fatal(err)
	}
	p.Connection = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	p.Channel = ch

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	p.Queue = q

	go func(p *AMQPClient) {
		for {
			select {
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

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func(p *AMQPClient) {
		for d := range msgs {
			p.incomingMessages <- &network.MessagePacket{Message: string(d.Body)}
		}
	}(p)

	p.stateUpdates <- true
}

func (p *AMQPClient) Disconnect() error {
	err := p.Connection.Close()
	if err != nil {
		return err
	}
	p.stateUpdates <- false
	return nil
}

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

func (p AMQPClient) Receive() chan *network.MessagePacket {
	return p.incomingMessages
}

func (p AMQPClient) State() chan bool {
	return p.stateUpdates
}
