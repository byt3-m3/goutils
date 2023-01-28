package ikafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

type (
	PublishMessageMockResponse struct {
		LinesWritten int
		Error        error
	}
)

type KafkaPublisherMock struct {
	mock.Mock

	PublishMessageMockResponse *PublishMessageMockResponse
}

func (i *KafkaPublisherMock) PublishMessage(ctx context.Context, msg *kafka.Message) (int, error) {
	return i.PublishMessageMockResponse.LinesWritten, i.PublishMessageMockResponse.Error
}
