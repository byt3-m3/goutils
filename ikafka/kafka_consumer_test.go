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

	//kafkaConsumerLive := func() IConsumer {
	//
	//	return ProvideKafkaConsumer(&KafkaConsumerConfig{
	//		Topic:      "test-topic",
	//		Brokers:    []string{"192.168.1.5:9092"},
	//		ConsumerID: "",
	//	})
	//}

	testKafkaTopic    = "test-topic"
	testKafkaMsgValue = []byte("This is data")
	kafkaConsumerErr  = errors.New("this is an error")

	kafkaConsumerMockNoErrs = func() IConsumer {
		return &KafkaConsumerMock{
			ConsumeMockResponse: &ConsumeMockResponse{
				Message: &kafka.Message{
					Topic:         testKafkaTopic,
					Partition:     0,
					Offset:        0,
					HighWaterMark: 0,
					Key:           nil,
					Value:         testKafkaMsgValue,
					Headers:       nil,
					Time:          time.Time{},
				},
				Error: nil,
			},
			ConsumeAsyncMockResponse: &ConsumeAsyncMockResponse{Error: nil},
		}
	}

	kafkaConsumerMockErrs = func() IConsumer {
		return &KafkaConsumerMock{
			ConsumeMockResponse: &ConsumeMockResponse{
				Message: nil,
				Error:   kafkaConsumerErr,
			},
			ConsumeAsyncMockResponse: &ConsumeAsyncMockResponse{Error: kafkaConsumerErr},
		}
	}
)

func TestKafkaConsumer_Consume(t *testing.T) {
	ctx := context.Background()

	t.Run("test when consumer is successful", func(t *testing.T) {
		c := dig.New()

		if err := c.Provide(kafkaConsumerMockNoErrs); err != nil {
			t.Fatal(err)
		}

		if err := c.Invoke(func(consumer IConsumer) {
			msg, err := consumer.Consume(ctx)
			assert.NoError(t, err)
			assert.Equal(t, testKafkaMsgValue, msg.Value)

		}); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("test when consumer failed", func(t *testing.T) {
		c := dig.New()

		if err := c.Provide(kafkaConsumerMockErrs); err != nil {
			t.Fatal(err)
		}

		if err := c.Invoke(func(consumer IConsumer) {
			msg, err := consumer.Consume(ctx)
			assert.Error(t, err)
			assert.Empty(t, msg)
			assert.Equal(t, kafkaConsumerErr, err)

		}); err != nil {
			t.Fatal(err)
		}

	})

}

func TestKafkaConsumer_ConsumeAsync(t *testing.T) {
	ctx := context.Background()

	t.Run("test when consumer is successful", func(t *testing.T) {
		c := dig.New()

		if err := c.Provide(kafkaConsumerMockNoErrs); err != nil {
			t.Fatal(err)
		}

		if err := c.Invoke(func(consumer IConsumer) {
			msgBus := make(chan *kafka.Message)

			err := consumer.ConsumeAsync(ctx, msgBus, time.Second*5)
			assert.NoError(t, err)

		}); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("test when consumer failed", func(t *testing.T) {
		c := dig.New()

		if err := c.Provide(kafkaConsumerMockErrs); err != nil {
			t.Fatal(err)
		}

		if err := c.Invoke(func(consumer IConsumer) {
			msgBus := make(chan *kafka.Message)

			err := consumer.ConsumeAsync(ctx, msgBus, time.Second*5)
			assert.Error(t, err)

		}); err != nil {
			t.Fatal(err)
		}

	})
}
