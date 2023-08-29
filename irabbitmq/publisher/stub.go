package publisher

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq/admin"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type StubRabbitMQPublisher struct {
	PublishStubReturn         func(ctx context.Context, input *PublishInput) error
	GetConnectionStubReturn   func() *amqp091.Connection
	ResetConnectionStubReturn func() error
	IsClosedStubReturn        func() bool
	WithAMQPUrlStubReturn     func(url string)
	WithVHostStubReturn       func(vhost string)
	WithLoggerStubReturn      func(logger *log.Logger)
	WithNoAuthStubReturn      func()
	WithPlainAuthStubReturn   func(username, password string)

	MustValidateStubReturn func()
}

func (s *StubRabbitMQPublisher) WithAMQPUrl(url string) RabbitMQPublisher {
	s.WithAMQPUrlStubReturn(url)
	return s

}

func (s *StubRabbitMQPublisher) WithVHost(vhost string) RabbitMQPublisher {
	s.WithVHostStubReturn(vhost)
	return s

}

func (s *StubRabbitMQPublisher) WithLogger(logger *log.Logger) RabbitMQPublisher {
	s.WithLoggerStubReturn(logger)
	return s
}

func (s *StubRabbitMQPublisher) WithNoAuth() RabbitMQPublisher {
	s.WithNoAuthStubReturn()
	return s
}

func (s *StubRabbitMQPublisher) WithPlainAuth(username, password string) RabbitMQPublisher {
	s.WithPlainAuthStubReturn(username, password)
	return s
}

func (s *StubRabbitMQPublisher) MustValidate() {
	s.MustValidateStubReturn()
}

type PublishStubReturn struct {
	Error error
}

func (s *StubRabbitMQPublisher) Publish(ctx context.Context, input *PublishInput) error {
	return s.PublishStubReturn(ctx, input)
}

type GetConnectionStubReturn struct {
	Conn *amqp091.Connection
}

func (s *StubRabbitMQPublisher) GetConnection() *amqp091.Connection {
	return s.GetConnectionStubReturn()
}

type GetAdminClientStubReturn struct {
	Client admin.Client
}

type ResetConnectionStubReturn struct {
	Error error
}

func (s *StubRabbitMQPublisher) ResetConnection() error {
	return s.ResetConnectionStubReturn()
}

type IsClosedStubReturn struct {
	IsClosed bool
}

func (s *StubRabbitMQPublisher) IsClosed() bool {
	return s.IsClosedStubReturn()
}
