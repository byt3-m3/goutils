package publisher

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"log/slog"
	"testing"
	"time"
)

var (
	stubPublisher = NewStubKafkaPublisher(&NewStubKafkaPublisherInput{
		PublishMessageStubReturn: func(ctx context.Context, msg *kafka.Message) (int, error) {
			return 0, nil
		},
		WithAuthStubReturn: func(authMechanism sasl.Mechanism) {

		},
		WithTopicStubReturn: func(topic string) {

		},
		WithBrokerStubReturn: func(broker string) {

		},
		WithKafkaConnStubReturn: func(conn *kafka.Conn) {

		},
		WithLoggerStubReturn: func(logger *slog.Logger) {

		},
	})
	kafkaPublisherErr = errors.New("this is an error")
	testKafkaMessage  = &kafka.Message{
		Topic:         "test-topic",
		Partition:     0,
		Offset:        0,
		HighWaterMark: 0,
		Key:           nil,
		Value:         nil,
		Headers:       nil,
		Time:          time.Time{},
	}
)

func TestKafkaPublisherMock_PublishMessage(t *testing.T) {

	t.Run("test when successful", func(t *testing.T) {

		_, _ = stubPublisher.PublishMessage(context.Background(), testKafkaMessage)

	})

}
