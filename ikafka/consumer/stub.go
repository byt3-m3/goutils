package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type StubKafkaConsumer struct {
	ConsumeAsyncStubReturn func(ctx context.Context, input *ConsumeAsyncInput) error
	ConsumeStubReturn      func(ctx context.Context) (*kafka.Message, error)
}

type ConsumeAsyncStubReturn struct {
	Error error
}

func (s *StubKafkaConsumer) ConsumeAsync(ctx context.Context, input *ConsumeAsyncInput) error {
	return s.ConsumeAsyncStubReturn(ctx, input)
}

type ConsumeStubReturn struct {
	Message *kafka.Message
	Error   error
}

func (s *StubKafkaConsumer) Consume(ctx context.Context) (*kafka.Message, error) {
	return s.ConsumeStubReturn(ctx)
}
