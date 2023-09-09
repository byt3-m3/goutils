package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type StubRabbitMQConsumer struct {
	ConsumeReturn               func(ctx context.Context, queue string) (<-chan amqp091.Delivery, error)
	GetConnectionReturn         func() irabbitmq.Connection
	IsClosedReturn              func() bool
	GetActiveChannelStubReturn  func() *amqp091.Channel
	WithAMQPUrlStubReturn       func(url string)
	WithConsumerIDStubReturn    func(id string)
	WithVHostStubReturn         func(vhost string)
	WithPlainAuthStubReturn     func(username, password string)
	WithPreFetchCountStubReturn func(count int)
	WithLoggerStubReturn        func(logger *log.Logger)
	ResetConnectionStubReturn   func()
}

func (s *StubRabbitMQConsumer) WithAMQPUrl(url string) RabbitMQConsumer {
	s.WithAMQPUrlStubReturn(url)
	return s
}

func (s *StubRabbitMQConsumer) WithConsumerID(id string) RabbitMQConsumer {
	s.WithConsumerIDStubReturn(id)
	return s
}

func (s *StubRabbitMQConsumer) WithVHost(vhost string) RabbitMQConsumer {
	s.WithVHostStubReturn(vhost)
	return s
}

func (s *StubRabbitMQConsumer) WithPlainAuth(username, password string) RabbitMQConsumer {
	s.WithPlainAuthStubReturn(username, password)
	return s
}

func (s *StubRabbitMQConsumer) WithPreFetchCount(count int) RabbitMQConsumer {
	s.WithPreFetchCountStubReturn(count)
	return s
}

func (s *StubRabbitMQConsumer) WithLogger(logger *log.Logger) RabbitMQConsumer {
	s.WithLoggerStubReturn(logger)
	return s
}

func (s *StubRabbitMQConsumer) ResetConnection() {
	s.ResetConnectionStubReturn()
}

func (s *StubRabbitMQConsumer) GetActiveChannel() *amqp091.Channel {
	return s.GetActiveChannelStubReturn()
}

func (s *StubRabbitMQConsumer) Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error) {
	return s.ConsumeReturn(ctx, queue)
}

func (s *StubRabbitMQConsumer) GetConnection() irabbitmq.Connection {
	return s.GetConnectionReturn()
}

func (s *StubRabbitMQConsumer) IsClosed() bool {
	return s.IsClosedReturn()
}
