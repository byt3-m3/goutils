package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type StubKafkaConsumer struct {
	ConsumeAsyncStubReturn func(ctx context.Context, msgBus chan *kafka.Message, tickerRate time.Duration) ConsumeAsyncStubReturn
	ConsumeStubReturn      func(ctx context.Context) ConsumeStubReturn
}

type ConsumeAsyncStubReturn struct {
	Error error
}

func (s *StubKafkaConsumer) ConsumeAsync(ctx context.Context, msgBus chan *kafka.Message, tickerRate time.Duration) error {
	return s.ConsumeAsyncStubReturn(ctx, msgBus, tickerRate).Error
}

type ConsumeStubReturn struct {
	Message *kafka.Message
	Error   error
}

func (s *StubKafkaConsumer) Consume(ctx context.Context) (*kafka.Message, error) {
	res := s.ConsumeStubReturn(ctx)
	return res.Message, res.Error
}
