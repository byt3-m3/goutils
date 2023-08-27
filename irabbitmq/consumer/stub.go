package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/irabbitmq/admin"
	"github.com/rabbitmq/amqp091-go"
)

type StubConsumer struct {
	ConsumeReturn              func(ctx context.Context, queue string) ConsumeStubReturn
	GetConnectionReturn        func() GetConnectionStubReturn
	GetAdminClientReturn       func() GetAdminClientStubReturn
	IsClosedReturn             func() IsClosedStubReturn
	GetActiveChannelStubReturn func() GetActiveChannelStubReturn
}

type GetActiveChannelStubReturn struct {
	Channel *amqp091.Channel
}

func (s *StubConsumer) GetActiveChannel() *amqp091.Channel {
	return s.GetActiveChannelStubReturn().Channel
}

type ConsumeStubReturn struct {
	DeliveryChan <-chan amqp091.Delivery
	Err          error
}

func (s *StubConsumer) Consume(ctx context.Context, queue string) (<-chan amqp091.Delivery, error) {
	res := s.ConsumeReturn(ctx, queue)
	return res.DeliveryChan, res.Err
}

type GetConnectionStubReturn struct {
	Conn *amqp091.Connection
}

func (s *StubConsumer) GetConnection() *amqp091.Connection {
	return s.GetConnectionReturn().Conn
}

type GetAdminClientStubReturn struct {
	Client admin.Client
}

func (s *StubConsumer) GetAdminClient() admin.Client {
	return s.GetAdminClientReturn().Client
}

type IsClosedStubReturn struct {
	Bool bool
}

func (s *StubConsumer) IsClosed() bool {
	return s.IsClosedReturn().Bool
}
