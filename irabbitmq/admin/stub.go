package admin

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

type StubClient struct {
	CreateQueueStubReturn    func(ctx context.Context, input *CreateQueueInput) (*amqp091.Queue, error)
	CreateExchangeStubReturn func(ctx context.Context, input *CreateExchangeInput) error
	BindQueueStubReturn      func(ctx context.Context, input *BindQueueInput) error
	DeleteQueueStubReturn    func(ctx context.Context, input *DeleteQueueInput) error
	DeleteExchangeStubReturn func(ctx context.Context, input *DeleteExchangeInput) error
	GetConnectionStubReturn  func() *amqp091.Connection
	WithAMQPUrlStubReturn    func(url string)
	WithVHostStubReturn      func(vhost string)
	WithPlainAuthStubReturn  func(username, password string)
	WithConnectionStubReturn func(conn *amqp091.Connection)
	ValidateClientStubReturn func(client *adminClient) bool
}

func (s *StubClient) WithAMQPUrl(url string) Client {
	s.WithAMQPUrlStubReturn(url)
	return s
}

type WithVHostStubReturn struct {
}

func (s *StubClient) WithVHost(vhost string) Client {
	s.WithVHostStubReturn(vhost)
	return s
}

func (s *StubClient) WithPlainAuth(username, password string) Client {
	s.WithPlainAuthStubReturn(username, password)
	return s
}

func (s *StubClient) WithConnection(conn *amqp091.Connection) Client {
	s.WithConnectionStubReturn(conn)
	return s
}

func (s *StubClient) ValidateClient(client *adminClient) bool {
	return s.ValidateClientStubReturn(client)
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

func (s *StubClient) GetConnection() *amqp091.Connection {
	return s.GetConnectionStubReturn()
}

func (s *StubClient) DeleteQueue(ctx context.Context, input *DeleteQueueInput) error {
	return s.DeleteQueueStubReturn(ctx, input)
}

func (s *StubClient) DeleteExchange(ctx context.Context, input *DeleteExchangeInput) error {
	return s.DeleteExchangeStubReturn(ctx, input)
}
