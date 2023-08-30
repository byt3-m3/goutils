package admin_client

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type QueueCreator interface {
	CreateQueue(ctx context.Context, input *CreateQueueInput) (*amqp091.Queue, error)
}

type QueueDeleter interface {
	DeleteQueue(ctx context.Context, input *DeleteQueueInput) error
}

type ExchangeCreator interface {
	CreateExchange(ctx context.Context, input *CreateExchangeInput) error
}

type ExchangeDeleter interface {
	DeleteExchange(ctx context.Context, input *DeleteExchangeInput) error
}

type QueueBinder interface {
	BindQueue(ctx context.Context, input *BindQueueInput) error
}

type ConnectionGetter interface {
	GetConnection() *amqp091.Connection
}

type ClientOptionSetter interface {
	WithAMQPUrl(url string) RabbitMQAdminClient

	WithVHost(vhost string) RabbitMQAdminClient

	WithPlainAuth(username, password string) RabbitMQAdminClient

	WithConnection(conn *amqp091.Connection) RabbitMQAdminClient

	WithLogger(logger *log.Logger) RabbitMQAdminClient
}

type ClientValidator interface {
	ValidateClient(client *adminClient) bool
}

type RabbitMQAdminClient interface {
	QueueCreator

	ExchangeCreator

	QueueBinder

	ConnectionGetter

	QueueDeleter

	ExchangeDeleter

	ClientOptionSetter

	ClientValidator
}
