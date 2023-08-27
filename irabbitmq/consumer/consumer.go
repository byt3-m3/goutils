package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq/admin"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type consumerOpt func(c *consumer)

var (
	WithAMQPUrl = func(url string) consumerOpt {
		return func(c *consumer) {
			c.amqpUrl = url
		}
	}

	WithConsumerID = func(id string) consumerOpt {
		return func(c *consumer) {
			c.consumerID = id
		}
	}

	WithVhost = func(vhost string) consumerOpt {
		return func(c *consumer) {
			c.vHost = vhost
		}
	}

	WithPlainAuth = func(username, password string) consumerOpt {
		return func(c *consumer) {
			c.amqpAuth = &amqp091.PlainAuth{
				Username: username,
				Password: password,
			}
		}
	}

	WithPreFetchCount = func(count int) consumerOpt {
		return func(c *consumer) {
			c.prefetchCount = count
		}
	}

	WithLogger = func(logger *log.Logger) consumerOpt {
		return func(c *consumer) {
			c.logger = logger
		}
	}
)

func validateConsumer(c *consumer) bool {
	if c.consumerID == "" {
		log.Println("consumer id not set, use WithConsumerID")
		return false
	}

	if c.amqpUrl == "" {
		log.Println("AMQPUrl not set. use WithAMPQURL")
		return false
	}

	if c.amqpUrl == "" {
		log.Println("AMQPUrl not set. use WithAMPQURL")
		return false
	}

	return true
}

type Consumer interface {
	Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error)
	GetConnection() *amqp091.Connection
	GetAdminClient() admin.Client
	IsClosed() bool
}

type consumer struct {
	amqpUrl       string
	consumerID    string
	vHost         string
	conn          *amqp091.Connection
	amqpAuth      amqp091.Authentication
	logger        *log.Logger
	prefetchCount int
	adminClient   admin.Client
	activeChannel *amqp091.Channel
}

func NewConsumer(opts ...consumerOpt) Consumer {
	c := &consumer{}

	for _, opt := range opts {
		opt(c)
	}

	if c.logger == nil {
		c.logger = log.Default()
	}

	if !validateConsumer(c) {
		c.logger.Fatalln("failed option validation")
	}

	if c.conn == nil {
		c.logger.Println("creating connection")
		c.ResetConnection()

	}

	if c.adminClient == nil {
		c.logger.Println("creating admin client")
		ac := admin.NewAdminClient(
			admin.WithConnection(c.conn),
			admin.WithAMQPUrl(c.amqpUrl),
		)

		c.adminClient = ac
	}

	return c

}

func (c *consumer) Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error) {

	ch, err := c.getChannel()

	if err := ch.Qos(c.prefetchCount, 0, false); err != nil {
		return nil, err

	}

	delivery, err := ch.Consume(queue, c.consumerID, false, false, false, false, nil)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	return delivery, nil

}

func (c *consumer) GetConnection() *amqp091.Connection {
	return c.conn
}

func (c *consumer) GetAdminClient() admin.Client {
	return c.adminClient
}

func (c *consumer) IsClosed() bool {
	return c.conn.IsClosed()
}

func (c *consumer) ResetConnection() {
	cConn, err := amqp091.DialConfig(c.amqpUrl, amqp091.Config{
		SASL:            []amqp091.Authentication{c.amqpAuth},
		Vhost:           c.vHost,
		ChannelMax:      0,
		FrameSize:       0,
		Heartbeat:       0,
		TLSClientConfig: nil,
		Properties:      nil,
		Locale:          "",
		Dial:            nil,
	})

	if err != nil {
		c.logger.Print(err)
	}

	c.conn = cConn
}

func (c *consumer) getChannel() (*amqp091.Channel, error) {
	if c.IsClosed() {
		c.ResetConnection()
	}

	ch, err := c.conn.Channel()
	if err != nil {
		c.logger.Println(err)
	}

	c.activeChannel = ch

	return c.activeChannel, nil
}

func (c *consumer) GetActiveChannel() *amqp091.Channel {
	return c.activeChannel
}
