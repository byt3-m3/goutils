package ikafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
	"time"
)

type (
	ConsumeAsyncMockResponse struct {
		Error error
	}

	ConsumeMockResponse struct {
		Message *kafka.Message
		Error   error
	}
)

type KafkaConsumerMock struct {
	mock.Mock
	ConsumeMockResponse      *ConsumeMockResponse
	ConsumeAsyncMockResponse *ConsumeAsyncMockResponse
}

func (k *KafkaConsumerMock) ConsumeAsync(ctx context.Context, msgBus chan *kafka.Message, tickerRate time.Duration) error {
	return k.ConsumeAsyncMockResponse.Error
}

func (k *KafkaConsumerMock) Consume(ctx context.Context) (*kafka.Message, error) {
	return k.ConsumeMockResponse.Message, k.ConsumeMockResponse.Error
}
