package protocol

import (
	"github.com/streadway/amqp"
	"github.com/superioz/gochat/internal/network"
)

const queueName string = "goqueue"

type AMQPClient struct {
	Connection       *amqp.Connection
	OutgoingMessages chan *network.MessagePacket
	IncomingMessages chan *network.MessagePacket
	StateUpdates     chan bool
	Channel          *amqp.Channel
	Queue            amqp.Queue
}

func (p *AMQPClient) Connect(ip string) error {
	conn, err := amqp.Dial(ip)
	if err != nil {
		return err
	}
	p.Connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	p.Channel = ch

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	p.Queue = q

	go func(p *AMQPClient) {
		for {
			select {
			case s := <-p.OutgoingMessages:
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
		return err
	}

	go func(p *AMQPClient) {
		for d := range msgs {
			p.IncomingMessages <- &network.MessagePacket{Message: string(d.Body)}
		}
	}(p)

	p.StateUpdates <- true
	return nil
}

func (p *AMQPClient) Disconnect() error {
	err := p.Connection.Close()
	if err != nil {
		return err
	}
	p.StateUpdates <- false
	return nil
}

func (p AMQPClient) Send() chan *network.MessagePacket {
	return p.OutgoingMessages
}

func (p AMQPClient) Receive() chan *network.MessagePacket {
	return p.IncomingMessages
}

func (p AMQPClient) State() chan bool {
	return p.StateUpdates
}
