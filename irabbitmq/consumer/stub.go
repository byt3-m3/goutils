package consumer

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type StubConsumer struct {
	ConsumeReturn               func(ctx context.Context, queue string) (<-chan amqp091.Delivery, error)
	GetConnectionReturn         func() *amqp091.Connection
	IsClosedReturn              func() bool
	GetActiveChannelStubReturn  func() *amqp091.Channel
	WithAMQPUrlStubReturn       func(url string) RabbitMQConsumer
	WithConsumerIDStubReturn    func(id string) RabbitMQConsumer
	WithVHostStubReturn         func(vhost string) RabbitMQConsumer
	WithPlainAuthStubReturn     func(username, password string) RabbitMQConsumer
	WithPreFetchCountStubReturn func(count int) RabbitMQConsumer
	WithLoggerStubReturn        func(logger *log.Logger) RabbitMQConsumer
	ResetConnectionStubReturn   func()
}

func (s *StubConsumer) WithAMQPUrl(url string) RabbitMQConsumer {
	s.WithAMQPUrlStubReturn(url)
	return s
}

func (s *StubConsumer) WithConsumerID(id string) RabbitMQConsumer {
	s.WithConsumerIDStubReturn(id)
	return s
}

func (s *StubConsumer) WithVHost(vhost string) RabbitMQConsumer {
	s.WithVHostStubReturn(vhost)
	return s
}

func (s *StubConsumer) WithPlainAuth(username, password string) RabbitMQConsumer {
	s.WithPlainAuthStubReturn(username, password)
	return s
}

func (s *StubConsumer) WithPreFetchCount(count int) RabbitMQConsumer {
	s.WithPreFetchCountStubReturn(count)
	return s
}

func (s *StubConsumer) WithLogger(logger *log.Logger) RabbitMQConsumer {
	s.WithLoggerStubReturn(logger)
	return s
}

func (s *StubConsumer) ResetConnection() {
	s.ResetConnectionStubReturn()
}

func (s *StubConsumer) GetActiveChannel() *amqp091.Channel {
	return s.GetActiveChannelStubReturn()
}

func (s *StubConsumer) Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error) {
	return s.ConsumeReturn(ctx, queue)
}

func (s *StubConsumer) GetConnection() *amqp091.Connection {
	return s.GetConnectionReturn()
}

func (s *StubConsumer) IsClosed() bool {
	return s.IsClosedReturn()
}
