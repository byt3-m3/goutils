package admin_client

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type StubClient struct {
	CreateQueueStubReturn    func(ctx context.Context, input *CreateQueueInput) (*amqp091.Queue, error)
	CreateExchangeStubReturn func(ctx context.Context, input *CreateExchangeInput) error
	BindQueueStubReturn      func(ctx context.Context, input *BindQueueInput) error
	DeleteQueueStubReturn    func(ctx context.Context, input *DeleteQueueInput) error
	DeleteExchangeStubReturn func(ctx context.Context, input *DeleteExchangeInput) error
	GetConnectionStubReturn  func() irabbitmq.Connection
	WithAMQPUrlStubReturn    func(url string)
	WithVHostStubReturn      func(vhost string)
	WithPlainAuthStubReturn  func(username, password string)
	WithConnectionStubReturn func(conn *amqp091.Connection)
	MustValidateStubReturn   func()
	WithLoggerStubReturn     func(logger *slog.Logger)
}

func (s *StubClient) WithLogger(logger *slog.Logger) RabbitMQAdminClient {
	s.WithLoggerStubReturn(logger)
	return s
}

func (s *StubClient) WithAMQPUrl(url string) RabbitMQAdminClient {
	s.WithAMQPUrlStubReturn(url)
	return s
}

type WithVHostStubReturn struct {
}

func (s *StubClient) WithVHost(vhost string) RabbitMQAdminClient {
	s.WithVHostStubReturn(vhost)
	return s
}

func (s *StubClient) WithPlainAuth(username, password string) RabbitMQAdminClient {
	s.WithPlainAuthStubReturn(username, password)
	return s
}

func (s *StubClient) WithConnection(conn *amqp091.Connection) RabbitMQAdminClient {
	s.WithConnectionStubReturn(conn)
	return s
}

func (s *StubClient) MustValidate() {
	s.MustValidateStubReturn()
}

func (s *StubClient) CreateQueue(ctx context.Context, input *CreateQueueInput) (*amqp091.Queue, error) {
	return s.CreateQueueStubReturn(ctx, input)
}

func (s *StubClient) CreateExchange(ctx context.Context, input *CreateExchangeInput) error {
	return s.CreateExchangeStubReturn(ctx, input)
}

func (s *StubClient) BindQueue(ctx context.Context, input *BindQueueInput) error {
	return s.BindQueueStubReturn(ctx, input)
}

func (s *StubClient) GetConnection() irabbitmq.Connection {
	return s.GetConnectionStubReturn()
}

func (s *StubClient) DeleteQueue(ctx context.Context, input *DeleteQueueInput) error {
	return s.DeleteQueueStubReturn(ctx, input)
}

func (s *StubClient) DeleteExchange(ctx context.Context, input *DeleteExchangeInput) error {
	return s.DeleteExchangeStubReturn(ctx, input)
}
