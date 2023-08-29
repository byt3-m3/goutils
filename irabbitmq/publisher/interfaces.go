package publisher

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMQPublisherConnectionRester interface {
	ResetConnection() error
}

type RabbitMQPublisherConnectionGetter interface {
	GetConnection() *amqp091.Connection
}

type RabbitMQPublisherConnectionChecker interface {
	IsClosed() bool
}

type RabbitMQPublisherOptionSetter interface {
	WithAMQPUrl(url string) RabbitMQPublisher

	WithVHost(vhost string) RabbitMQPublisher
	WithLogger(log *log.Logger) RabbitMQPublisher
	WithNoAuth() RabbitMQPublisher

	WithPlainAuth(username, password string) RabbitMQPublisher
}

type RabbitMQPublisherValidator interface {
	MustValidate()
}

type RabbitMQPublisher interface {
	RabbitMQPublisherOptionSetter
	RabbitMQPublisherValidator
	RabbitMQPublisherConnectionRester
	RabbitMQPublisherConnectionGetter
	RabbitMQPublisherConnectionChecker
	Publish(ctx context.Context, input *PublishInput) error
}
