package ikafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"testing"
	"time"
)

var (
	kafkaPublisherMockNoErrs = func() IPublisher {
		return &KafkaPublisherMock{
			PublishMessageMockResponse: &PublishMessageMockResponse{
				LinesWritten: 30,
				Error:        nil,
			},
		}
	}

	kafkaPublisherMockErrs = func() IPublisher {
		return &KafkaPublisherMock{
			PublishMessageMockResponse: &PublishMessageMockResponse{
				LinesWritten: 0,
				Error:        kafkaPublisherErr,
			},
		}
	}
	kafkaPublisherErr = errors.New("this is an error")
	testKafkaMessage  = &kafka.Message{
		Topic:         testKafkaTopic,
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
		c := dig.New()
		if err := c.Provide(kafkaPublisherMockNoErrs); err != nil {
			t.Fatal(err)
		}

		if err := c.Invoke(func(publisher IPublisher) {
			ctx := context.Background()
			lineWritten, err := publisher.PublishMessage(ctx, testKafkaMessage)
			assert.NoError(t, err)

			assert.Equal(t, 30, lineWritten)

		}); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("test when failed", func(t *testing.T) {
		c := dig.New()
		if err := c.Provide(kafkaPublisherMockErrs); err != nil {
			t.Fatal(err)
		}

		if err := c.Invoke(func(publisher IPublisher) {
			ctx := context.Background()
			lineWritten, err := publisher.PublishMessage(ctx, testKafkaMessage)
			assert.Error(t, err)

			assert.Equal(t, 0, lineWritten)

		}); err != nil {
			t.Fatal(err)
		}

	})

}
