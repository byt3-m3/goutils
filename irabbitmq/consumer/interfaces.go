package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq/connection_handler"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumerValidator interface {
	MustValidate()
}

type RabbitMQConsumer interface {
	Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error)
	GetConnectionHandler() connection_handler.ConnectionHandler
}
