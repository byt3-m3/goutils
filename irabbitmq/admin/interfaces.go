package admin

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
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
	WithAMQPUrl(url string) Client

	WithVHost(vhost string) Client

	WithPlainAuth(username, password string) Client

	WithConnection(conn *amqp091.Connection) Client
}

type ClientValidator interface {
	ValidateClient(client *adminClient) bool
}

type Client interface {
	QueueCreator

	ExchangeCreator

	QueueBinder

	ConnectionGetter

	QueueDeleter

	ExchangeDeleter

	ClientOptionSetter

	ClientValidator
}