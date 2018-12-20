package protocol

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/superioz/gochat/internal/env"
	"github.com/superioz/gochat/internal/logs"
	"github.com/superioz/gochat/internal/network"
	"log"
)

const queueName string = "goqueue"

// represents an amqp client
type AMQPClient struct {
	UUID             uuid.UUID
	Nick             string
	Connection       *amqp.Connection
	outgoingMessages chan *network.MessagePacket
	incomingMessages chan *network.MessagePacket
	stateUpdates     chan bool
	Channel          *amqp.Channel
	Queue            amqp.Queue
	Logger           logs.ChatLogger
}

func NewAMQPClient(nick string) AMQPClient {
	return AMQPClient{UUID: uuid.NewV4(), Nick: nick, outgoingMessages: make(chan *network.MessagePacket),
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
	fmt.Println("Connected to amqp server @" + ip + ".")

	// start logging
	if env.IsLoggingEnabled() {
		go func(p *AMQPClient) {
			p.Logger, err = logs.CreateAndConnect(env.GetLoggingCredentials())

			if err != nil {
				fmt.Println("Couldn't connect to logging service! No logs will be stored..")
			}
		}(p)
	}

	// creates a channel to amqp server
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	p.Channel = ch

	// declares an exchange with name `chat` and type `fanout`
	// sends to every queue bound to this exchange
	err = ch.ExchangeDeclare("chat", "fanout", true,
		false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// declares a new queue with name=`queueName`
	q, err := ch.QueueDeclare(queueName+"_"+p.UUID.String(), false, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	p.Queue = q

	// binds the queue to the exchange
	err = ch.QueueBind(q.Name, "", "chat", false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// handles outgoing messages
	go func(p *AMQPClient) {
		for {
			select {
			// publishes a new message to the amqp server
			// if the channel received a new message
			case s := <-p.outgoingMessages:
				err = ch.Publish("chat", "", false, false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(s.Message),
					})
			case m := <-p.incomingMessages:
				fmt.Println(m.Message)

				// log the message
				if env.IsLoggingEnabled() {
					go func() {
						user, message := m.UserAndMessage()
						err := p.Logger.AddEntry(user, message)

						if p.Logger.Connected && err != nil {
							fmt.Println("Couldn't sent log!", err)
						}
					}()
				}
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

// returns the current nickname
func (p AMQPClient) Nickname() string {
	return p.Nick
}
