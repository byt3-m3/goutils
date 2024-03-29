package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq/connection_handler"
	"github.com/byt3-m3/goutils/logging"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

type consumer struct {
	consumerID    string
	logger        *slog.Logger
	prefetchCount int
	connHandler   connection_handler.ConnectionHandler
}

type NewOpt func(c *consumer)

func WithConsumerID(consumerID string) NewOpt {
	return func(c *consumer) {
		c.consumerID = consumerID
	}
}

func WithPrefetchCount(count int) NewOpt {
	return func(c *consumer) {
		c.prefetchCount = count
	}
}

func WithLogger(logger *slog.Logger) NewOpt {
	return func(c *consumer) {
		c.logger = logger
	}
}

func WithConnectionHandler(handler connection_handler.ConnectionHandler) NewOpt {

	return func(c *consumer) {
		c.connHandler = handler
	}
}

func New(opts ...NewOpt) RabbitMQConsumer {
	c := &consumer{}

	for _, opt := range opts {
		opt(c)
	}

	c.MustValidate()

	return c

}

func (c *consumer) MustValidate() {

	if c.consumerID == "" {
		c.consumerID = primitive.NewObjectID().Hex()
	}

	if c.logger == nil {
		c.logger = logging.NewJSONLogger(slog.LevelInfo, false)
	}

	if c.connHandler == nil {
		panic("connection handler not set, use WithConnectionHandler")
	}

}

func (c *consumer) Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error) {
	c.MustValidate()

	ch := c.connHandler.MustGetChannel()

	if err := ch.Qos(c.prefetchCount, 0, false); err != nil {
		return nil, err

	}

	delivery, err := ch.Consume(queue, c.consumerID, false, false, false, false, nil)
	if err != nil {
		c.logger.Error("error consuming message",
			slog.Any("error", err),
		)
		return nil, err
	}

	return delivery, nil

}

func (c *consumer) GetConnectionHandler() connection_handler.ConnectionHandler {
	return c.connHandler
}
