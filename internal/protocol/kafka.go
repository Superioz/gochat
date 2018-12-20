package protocol

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/satori/go.uuid"
	"github.com/superioz/gochat/internal/env"
	"github.com/superioz/gochat/internal/logs"
	"github.com/superioz/gochat/internal/network"
	"log"
)

const topicName string = "gochat"

// represents a kafka client
type KafkaClient struct {
	UUID             uuid.UUID
	Nick             string
	outgoingMessages chan *network.MessagePacket
	incomingMessages chan *network.MessagePacket
	stateUpdates     chan bool
	Logger           logs.ChatLogger

	Consumer          sarama.Consumer
	PartitionConsumer sarama.PartitionConsumer
	Producer          sarama.AsyncProducer
}

func NewKafkaClient(nick string) KafkaClient {
	return KafkaClient{UUID: uuid.NewV4(), Nick: nick, outgoingMessages: make(chan *network.MessagePacket),
		incomingMessages: make(chan *network.MessagePacket), stateUpdates: make(chan bool)}
}

// connects the client to the amqp server
// uses the `ip` to connect to the server
func (c *KafkaClient) Connect(ip string) {
	go c.startConsumer(ip)
	go c.startProducer(ip)
}

// initialises the producer
func (c *KafkaClient) startProducer(ip string) {
	producer, err := sarama.NewAsyncProducer([]string{ip}, nil)
	if err != nil {
		panic(err)
	}
	c.Producer = producer
	fmt.Println("Connected producer to kafka node @" + ip + ".")

	select {
	case c.stateUpdates <- true:
	default:
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case msg := <-c.outgoingMessages:
			producer.Input() <- &sarama.ProducerMessage{Topic: topicName, Key: nil, Value: sarama.StringEncoder(msg.Message)}
			break
		case err := <-producer.Errors():
			log.Println("Failed to produce message", err)
			break
		}
	}
}

// initialises the consumer
func (c *KafkaClient) startConsumer(ip string) {
	consumer, err := sarama.NewConsumer([]string{ip}, nil)
	if err != nil {
		panic(err)
	}
	c.Consumer = consumer
	fmt.Println("Connected consumer to kafka node @" + ip + ".")

	// start logging
	if env.IsLoggingEnabled() {
		go func(c *KafkaClient) {
			c.Logger, err = logs.CreateAndConnect(env.GetLoggingCredentials())

			if err != nil {
				fmt.Println("Couldn't connect to logging service! No logs will be stored..")
			}
		}(c)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topicName, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	c.PartitionConsumer = partitionConsumer

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			m := &network.MessagePacket{Message: string(msg.Value)}

			fmt.Println(m.Message)

			// log the message
			if env.IsLoggingEnabled() {
				go func() {
					user, message := m.UserAndMessage()
					err := c.Logger.AddEntry(user, message)

					if c.Logger.Connected && err != nil {
						fmt.Println("Couldn't sent log!", err)
					}
				}()
			}
			break
		}
	}
}

// closes the kafka consumer and producer
// returns an error if not successful
func (c *KafkaClient) Disconnect() error {
	if err := c.Consumer.Close(); err != nil {
		return err
	}
	if err := c.Producer.Close(); err != nil {
		return err
	}
	c.stateUpdates <- false

	return nil
}

// sends a message packet to the server
func (c KafkaClient) Send(packet network.MessagePacket) {
	select {
	case c.outgoingMessages <- &packet:
		// successful
		break
	default:
		// not successful
		break
	}
}

// returns receive channel
func (c KafkaClient) Receive() chan *network.MessagePacket {
	return c.incomingMessages
}

// returns the current connection state channel
func (c KafkaClient) State() chan bool {
	return c.stateUpdates
}

// returns the current nickname
func (c KafkaClient) Nickname() string {
	return c.Nick
}
