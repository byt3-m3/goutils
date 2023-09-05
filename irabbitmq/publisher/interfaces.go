package publisher

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq"
	log "github.com/sirupsen/logrus"
)

type RabbitMQPublisherConnectionRester interface {
	ResetConnection() error
}

type RabbitMQPublisherConnectionGetter interface {
	GetConnection() irabbitmq.Connection
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
