package consumer

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"log"
	"testing"
	"time"
)

var (
	kafkaConsumerLive = func() Consumer {

		return ProvideKafkaConsumer(&options{
			Topic:      "test-topic",
			Brokers:    []string{"192.168.1.5:9092"},
			ConsumerID: "",
		})
	}

	testKafkaTopic    = "test-topic"
	testKafkaMsgValue = []byte("This is data")
	kafkaConsumerErr  = errors.New("this is an error")

	kafkaConsumerMockNoErrs = func() Consumer {
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

	kafkaConsumerMockErrs = func() Consumer {
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
		consumer := NewConsumer(
			WithBrokers([]string{"192.168.1.60:9094"}),
			WithPlainAuth(plain.Mechanism{
				Username: "cbaxter",
				Password: "kafka1",
			}),
			WithTopic("test-topic"),
			WithConsumerID("test-consumer-id-1"),
		)

		msg, err := consumer.Consume(context.Background())
		if err != nil {
			log.Println(err)
		}

		log.Println(msg)

	})

	t.Run("test when consumer failed", func(t *testing.T) {
		consumer := NewConsumer(
			WithBrokers([]string{"127.0.0.1:9092"}),
			WithConsumerID("test_consumer"),
			WithTopic("test-topic"),
		)

		msg, err := consumer.Consume(ctx)
		fmt.Println(msg, err)
		//assert.Error(t, err)
		//assert.Empty(t, msg)
		//assert.Equal(t, kafkaConsumerErr, err)

	})

}

func TestKafkaConsumer_ConsumeAsync(t *testing.T) {
	ctx := context.Background()

	t.Run("test when consumer is successful", func(t *testing.T) {
		c := dig.New()

		if err := c.Provide(kafkaConsumerMockNoErrs); err != nil {
			t.Fatal(err)
		}

		if err := c.Invoke(func(consumer Consumer) {
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

		if err := c.Invoke(func(consumer Consumer) {
			msgBus := make(chan *kafka.Message)

			err := consumer.ConsumeAsync(ctx, msgBus, time.Second*5)
			assert.Error(t, err)

		}); err != nil {
			t.Fatal(err)
		}

	})
}
