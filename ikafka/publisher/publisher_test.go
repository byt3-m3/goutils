package publisher

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"testing"
	"time"
)

var (
	kafkaPublisherMockNoErrs = func() Publisher {
		return &KafkaPublisherMock{
			PublishMessageMockResponse: &PublishMessageMockResponse{
				LinesWritten: 30,
				Error:        nil,
			},
		}
	}

	kafkaPublisherMockErrs = func() Publisher {
		return &KafkaPublisherMock{
			PublishMessageMockResponse: &PublishMessageMockResponse{
				LinesWritten: 0,
				Error:        kafkaPublisherErr,
			},
		}
	}
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
		p, _ := NewPublisher(
			WithBroker("192.168.1.60:9094"),
			WithTopic("test-topic"),
			WithPlainAuth(plain.Mechanism{
				Username: "cbaxter",
				Password: "kafka1",
			}),
		)

		_, _ = p.PublishMessage(context.Background(), testKafkaMessage)

	})

	t.Run("test when failed", func(t *testing.T) {
		c := dig.New()
		if err := c.Provide(kafkaPublisherMockErrs); err != nil {
			t.Fatal(err)
		}

		if err := c.Invoke(func(publisher Publisher) {
			ctx := context.Background()
			lineWritten, err := publisher.PublishMessage(ctx, testKafkaMessage)
			assert.Error(t, err)

			assert.Equal(t, 0, lineWritten)

		}); err != nil {
			t.Fatal(err)
		}

	})

}
