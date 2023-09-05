package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq"
	"github.com/byt3-m3/goutils/logging"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type consumer struct {
	amqpUrl       string
	consumerID    string
	vHost         string
	conn          *amqp091.Connection
	amqpAuth      amqp091.Authentication
	logger        *log.Logger
	prefetchCount int
	activeChannel *amqp091.Channel
}

func New() RabbitMQConsumer {
	c := &consumer{}

	return c

}

func (c *consumer) MustValidate() {

	if c.consumerID == "" {
		panic("consumerID not set")
	}

	if c.amqpUrl == "" {
		panic("AMQPUrl not set")

	}

	if c.logger == nil {
		c.logger = logging.NewLogger()
	}

	if c.conn == nil {
		c.ResetConnection()
	}

}

func (c *consumer) WithAMQPUrl(url string) RabbitMQConsumer {
	c.amqpUrl = url
	return c
}

func (c *consumer) WithConsumerID(id string) RabbitMQConsumer {
	c.consumerID = id
	return c
}

func (c *consumer) WithVHost(vhost string) RabbitMQConsumer {
	c.vHost = vhost
	return c
}

func (c *consumer) WithPlainAuth(username, password string) RabbitMQConsumer {
	c.amqpAuth = &amqp091.PlainAuth{
		Username: username,
		Password: password,
	}
	return c
}

func (c *consumer) WithPreFetchCount(count int) RabbitMQConsumer {
	c.prefetchCount = count
	return c
}

func (c *consumer) WithLogger(logger *log.Logger) RabbitMQConsumer {
	c.logger = logger
	return c
}

func (c *consumer) Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error) {
	c.MustValidate()

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

func (c *consumer) GetConnection() irabbitmq.Connection {
	return c.conn
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
