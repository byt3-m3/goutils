package admin

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

type StubClient struct {
	CreateQueueStubReturn    func(ctx context.Context, input *CreateQueueInput) CreateQueueStubReturn
	CreateExchangeStubReturn func(ctx context.Context, input *CreateExchangeInput) CreateExchangeStubReturn
	BindQueueStubReturn      func(ctx context.Context, input *BindQueueInput) BindQueueStubReturn
	DeleteQueueStubReturn    func(ctx context.Context, input *DeleteQueueInput) DeleteQueueStubReturn
	DeleteExchangeStubReturn func(ctx context.Context, input *DeleteExchangeInput) DeleteExchangeStubReturn
	GetConnectionStubReturn  func() GetConnectionStubReturn
}

func (s *StubClient) WithAMQPUrl(url string) *adminClient {
	//TODO implement me
	panic("implement me")
}

func (s *StubClient) WithVHost(vhost string) *adminClient {
	//TODO implement me
	panic("implement me")
}

func (s *StubClient) WithPlainAuth(username, password string) *adminClient {
	//TODO implement me
	panic("implement me")
}

func (s *StubClient) WithConnection(conn *amqp091.Connection) *adminClient {
	//TODO implement me
	panic("implement me")
}

func (s *StubClient) ValidateClient(client *adminClient) bool {
	//TODO implement me
	panic("implement me")
}

type CreateQueueStubReturn struct {
	Queue *amqp091.Queue
	Err   error
}

func (s *StubClient) CreateQueue(ctx context.Context, input *CreateQueueInput) (*amqp091.Queue, error) {
	res := s.CreateQueueStubReturn(ctx, input)
	return res.Queue, res.Err
}

type CreateExchangeStubReturn struct {
	Error error
}

func (s *StubClient) CreateExchange(ctx context.Context, input *CreateExchangeInput) error {
	return s.CreateExchangeStubReturn(ctx, input).Error
}

type BindQueueStubReturn struct {
	Error error
}

func (s *StubClient) BindQueue(ctx context.Context, input *BindQueueInput) error {
	return s.BindQueueStubReturn(ctx, input).Error
}

type GetConnectionStubReturn struct {
	Conn *amqp091.Connection
}

func (s *StubClient) GetConnection() *amqp091.Connection {
	return s.GetConnectionStubReturn().Conn
}

type DeleteQueueStubReturn struct {
	Error error
}

func (s *StubClient) DeleteQueue(ctx context.Context, input *DeleteQueueInput) error {
	return s.DeleteQueueStubReturn(ctx, input).Error
}

type DeleteExchangeStubReturn struct {
	Error error
}

func (s *StubClient) DeleteExchange(ctx context.Context, input *DeleteExchangeInput) error {
	return s.DeleteExchangeStubReturn(ctx, input).Error
}
