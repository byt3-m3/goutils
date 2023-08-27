package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type StubKafkaPublisher struct {
	PublishMessageStubReturn func(ctx context.Context, msg *kafka.Message) PublishMessageStubReturn
}

type PublishMessageStubReturn struct {
	Count int
	Error error
}

func (s StubKafkaPublisher) PublishMessage(ctx context.Context, msg *kafka.Message) (int, error) {
	res := s.PublishMessageStubReturn(ctx, msg)
	return res.Count, res.Error
}
