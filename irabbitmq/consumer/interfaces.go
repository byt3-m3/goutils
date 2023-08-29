package consumer

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMQConsumerConnectionGetter interface {
	GetConnection() *amqp091.Connection
}

type RabbitMQConsumerActiveChannelGetter interface {
	GetActiveChannel() *amqp091.Channel
}

type RabbitMQConsumerConnectionChecker interface {
	IsClosed() bool
}

type RabbitMQConsumerConnectionRester interface {
	ResetConnection()
}

type RabbitMQConsumerOptionSetter interface {
	WithAMQPUrl(url string) RabbitMQConsumer
	WithConsumerID(id string) RabbitMQConsumer
	WithVHost(vhost string) RabbitMQConsumer
	WithPlainAuth(username, password string) RabbitMQConsumer
	WithPreFetchCount(count int) RabbitMQConsumer
	WithLogger(logger *log.Logger) RabbitMQConsumer
}

type RabbitMQConsumerValidator interface {
	MustValidate()
}

type RabbitMQConsumer interface {
	RabbitMQConsumerOptionSetter
	RabbitMQConsumerConnectionGetter
	RabbitMQConsumerConnectionChecker
	RabbitMQConsumerActiveChannelGetter
	RabbitMQConsumerConnectionRester
	Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error)
}
