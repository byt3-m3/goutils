package publisher

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq/admin"
	"github.com/rabbitmq/amqp091-go"
)

type StubPublisher struct {
	PublishStubReturn         func(ctx context.Context, input *PublishInput) PublishStubReturn
	GetConnectionStubReturn   func() GetConnectionStubReturn
	GetAdminClientStubReturn  func() GetAdminClientStubReturn
	ResetConnectionStubReturn func() ResetConnectionStubReturn
	IsClosedStubReturn        func() IsClosedStubReturn
}

type PublishStubReturn struct {
	Error error
}

func (s *StubPublisher) Publish(ctx context.Context, input *PublishInput) error {
	return s.PublishStubReturn(ctx, input).Error
}

type GetConnectionStubReturn struct {
	Conn *amqp091.Connection
}

func (s *StubPublisher) GetConnection() *amqp091.Connection {
	return s.GetConnectionStubReturn().Conn
}

type GetAdminClientStubReturn struct {
	Client admin.Client
}

func (s *StubPublisher) GetAdminClient() admin.Client {
	return s.GetAdminClientStubReturn().Client
}

type ResetConnectionStubReturn struct {
	Error error
}

func (s *StubPublisher) ResetConnection() error {
	return s.ResetConnectionStubReturn().Error
}

type IsClosedStubReturn struct {
	IsClosed bool
}

func (s *StubPublisher) IsClosed() bool {
	return s.IsClosedStubReturn().IsClosed
}
