package admin_client

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type adminClient struct {
	amqpUrl  string
	conn     *amqp091.Connection
	amqpAuth amqp091.Authentication
	vHost    string
	logger   *slog.Logger
}

func New() RabbitMQAdminClient {
	c := &adminClient{}

	return c

}

func (c *adminClient) MustValidate() {

	if c.logger == nil {
		c.logger = slog.Default()

	}

	if c.amqpUrl == "" {
		panic("amqpUrl not set, use WithAMQPUrl")

	}

	if c.conn == nil {
		conn, err := amqp091.DialConfig(c.amqpUrl, amqp091.Config{
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
			c.logger.Error("unable to dial rabbitmq server",
				slog.Any("error", err),
				slog.String("url", c.amqpUrl),
				slog.String("vhost", c.vHost),
			)
			panic(err)
		}

		c.conn = conn
	}

}

func (c *adminClient) WithLogger(logger *slog.Logger) RabbitMQAdminClient {
	c.logger = logger
	return c
}

func (c *adminClient) WithAMQPUrl(url string) RabbitMQAdminClient {
	c.amqpUrl = url

	return c
}

func (c *adminClient) WithVHost(vhost string) RabbitMQAdminClient {
	c.vHost = vhost
	return c
}
func (c *adminClient) WithPlainAuth(username, password string) RabbitMQAdminClient {
	c.amqpAuth = &amqp091.PlainAuth{
		Username: username,
		Password: password,
	}
	return c
}

func (c *adminClient) WithConnection(conn *amqp091.Connection) RabbitMQAdminClient {
	c.conn = conn

	return c
}

type CreateQueueInput struct {
	QueName       string
	IsDurable     bool
	IsExclusive   bool
	CanAutoDelete bool
	CanNoWait     bool
}

func (c *adminClient) CreateQueue(ctx context.Context, input *CreateQueueInput) (*amqp091.Queue, error) {

	c.MustValidate()
	ch, err := c.getChannel()
	if err != nil {
		c.logger.Error("unable to get channel",
			slog.Any("error", err),
			slog.String("queue", input.QueName),
		)
		return nil, err
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(input.QueName, input.IsDurable, input.CanAutoDelete, input.IsExclusive, input.CanNoWait, nil)
	if err != nil {
		c.logger.Error("unable to declare queue",
			slog.Any("error", err),
			slog.String("queue", input.QueName),
		)
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

	c.MustValidate()
	ch, err := c.getChannel()
	if err != nil {
		c.logger.Error("unable to get channel",
			slog.Any("error", err),
			slog.String("queue", input.QueName),
		)
		return err
	}

	defer ch.Close()

	if err := ch.QueueBind(input.QueName, input.Key, input.Exchange, input.CanNoWait, input.Args); err != nil {
		amqpErr := err.(*amqp091.Error)
		c.logger.Error("unable to bind queue",
			slog.Any("error", err),
			slog.String("queue", input.QueName),
		)
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
	c.MustValidate()
	ch, err := c.getChannel()
	if err != nil {
		c.logger.Error("unable to get channel",
			slog.Any("error", err),
			slog.String("queue", input.ExchangeName),
		)
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(input.ExchangeName, string(input.ExchangeType), input.IsDurable, input.CanAutoDelete, input.IsInternal, input.CanNoWait, nil)
	if err != nil {
		amqpErr := err.(amqp091.Error)
		c.logger.Error("unable to declare exchange",
			slog.Any("error", amqpErr),
			slog.String("exchange_name", input.ExchangeName),
			slog.Any("exchange_type", input.ExchangeType),
		)
		return amqpErr
	}

	return nil

}

func (c *adminClient) getChannel() (*amqp091.Channel, error) {

	if c.conn.IsClosed() || c.conn == nil {
		c.mustSetConnection()
	}

	return c.conn.Channel()

}

func (c *adminClient) GetConnection() irabbitmq.Connection {

	c.MustValidate()
	return c.conn
}

type DeleteQueueInput struct {
	Name     string
	IfUnUsed bool
	IfEmpty  bool
	NoWait   bool
}

func (c *adminClient) DeleteQueue(ctx context.Context, input *DeleteQueueInput) error {
	c.MustValidate()
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDelete(input.Name, input.IfUnUsed, input.IfEmpty, input.NoWait)
	if err != nil {
		c.logger.Error("error deleting queue",
			slog.Any("error", err),
		)
		return err
	}

	return err
}

type DeleteExchangeInput struct {
	Name     string
	IfUnUsed bool
	NoWait   bool
}

func (c *adminClient) DeleteExchange(ctx context.Context, input *DeleteExchangeInput) error {

	c.MustValidate()
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	if err := ch.ExchangeDelete(input.Name, input.IfUnUsed, input.NoWait); err != nil {
		c.logger.Error("error deleting exchange",
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

func (c *adminClient) mustSetConnection() {

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
		panic(err)
	}

	c.conn = cConn
}
