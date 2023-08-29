package admin

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type adminClient struct {
	amqpUrl  string
	conn     *amqp091.Connection
	amqpAuth amqp091.Authentication
	vHost    string
}

func New() Client {
	c := &adminClient{}

	return c

}

func (c *adminClient) ValidateClient(client *adminClient) bool {

	switch {
	case client.amqpUrl == "":
		log.Println("amqpUrl not set, use WithAMQPUrl")
		return false

	case client.conn == nil:
		conn, err := amqp091.DialConfig(client.amqpUrl, amqp091.Config{
			SASL:            []amqp091.Authentication{client.amqpAuth},
			Vhost:           client.vHost,
			ChannelMax:      0,
			FrameSize:       0,
			Heartbeat:       0,
			TLSClientConfig: nil,
			Properties:      nil,
			Locale:          "",
			Dial:            nil,
		})

		if err != nil {
			log.Print(err)
		}

		client.conn = conn

	}

	return true
}

func (c *adminClient) WithAMQPUrl(url string) Client {
	c.amqpUrl = url

	return c
}

func (c *adminClient) WithVHost(vhost string) Client {
	c.vHost = vhost
	return c
}
func (c *adminClient) WithPlainAuth(username, password string) Client {
	c.amqpAuth = &amqp091.PlainAuth{
		Username: username,
		Password: password,
	}
	return c
}

func (c *adminClient) WithConnection(conn *amqp091.Connection) Client {
	c.conn = conn

	return nil
}

type CreateQueueInput struct {
	QueName       string
	IsDurable     bool
	IsExclusive   bool
	CanAutoDelete bool
	CanNoWait     bool
}

func (c *adminClient) CreateQueue(ctx context.Context, input *CreateQueueInput) (*amqp091.Queue, error) {

	ch, err := c.getChannel()
	if err != nil {
		log.Print(err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(input.QueName, input.IsDurable, input.CanAutoDelete, input.IsExclusive, input.CanNoWait, nil)
	if err != nil {
		return nil, err
	}

	return &q, err

}

type BindQueueInput struct {
	QueName   string
	Key       string
	Exchange  string
	CanNoWait bool
	Args      amqp091.Table
}

func (c *adminClient) BindQueue(ctx context.Context, input *BindQueueInput) error {
	ch, err := c.getChannel()
	if err != nil {
		log.Print(err)
	}

	defer ch.Close()

	if err := ch.QueueBind(input.QueName, input.Key, input.Exchange, input.CanNoWait, input.Args); err != nil {
		amqpErr := err.(*amqp091.Error)

		log.Println(err)
		return amqpErr

	}

	return nil

}

type CreateExchangeInput struct {
	ExchangeName  string
	ExchangeType  AMQPExchangeType
	IsDurable     bool
	IsInternal    bool
	CanAutoDelete bool
	CanNoWait     bool
}

func (c *adminClient) CreateExchange(ctx context.Context, input *CreateExchangeInput) error {
	ch, err := c.getChannel()
	if err != nil {
		log.Print(err)
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(input.ExchangeName, string(input.ExchangeType), input.IsDurable, input.CanAutoDelete, input.IsInternal, input.CanNoWait, nil)
	if err != nil {
		amqpErr := err.(amqp091.Error)
		log.Println(amqpErr)
		return err
	}

	return nil

}

func (c *adminClient) getChannel() (*amqp091.Channel, error) {
	if c.conn.IsClosed() {
		c.setConnection()
	}

	return c.conn.Channel()

}

func (c *adminClient) GetConnection() *amqp091.Connection {
	return c.conn
}

type DeleteQueueInput struct {
	Name     string
	IfUnUsed bool
	IfEmpty  bool
	NoWait   bool
}

func (c *adminClient) DeleteQueue(ctx context.Context, input *DeleteQueueInput) error {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = ch.QueueDelete(input.Name, input.IfUnUsed, input.IfEmpty, input.NoWait)

	return err
}

type DeleteExchangeInput struct {
	Name     string
	IfUnUsed bool
	NoWait   bool
}

func (c *adminClient) DeleteExchange(ctx context.Context, input *DeleteExchangeInput) error {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Println(err)
		return err
	}

	return ch.ExchangeDelete(input.Name, input.IfUnUsed, input.NoWait)

}

func (c *adminClient) setConnection() {
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
		log.Println(err)
	}

	c.conn = cConn
}
