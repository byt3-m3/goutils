package consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

var (
	consumerStub = NewStub(NewStubInput{
		ConsumeAsyncStubReturn: func(ctx context.Context, msgBus chan *kafka.Message, tickerRate time.Duration) error {
			return nil
		},
		ConsumeStubReturn: func(ctx context.Context) (*kafka.Message, error) {
			return &kafka.Message{}, nil
		},
		WithReaderStubReturn: func(reader *kafka.Reader) {

		},
		WithTopicStubReturn: func(topic string) {

		},
		WithBrokersStubReturn: func(brokers []string) {

		},
		WithConsumerIDStubReturn: func(consumerID string) {

		},
		WithAuthStubReturn: func(authMechanism sasl.Mechanism) {

		},
	})

	consumerStubPanic = NewStub(NewStubInput{
		ConsumeAsyncStubReturn: func(ctx context.Context, msgBus chan *kafka.Message, tickerRate time.Duration) error {
			panic("oh no panic")

			return nil
		},
		ConsumeStubReturn: func(ctx context.Context) (*kafka.Message, error) {
			panic("oh no panic")

			return &kafka.Message{}, nil
		},
		WithReaderStubReturn: func(reader *kafka.Reader) {

		},
		WithTopicStubReturn: func(topic string) {

		},
		WithBrokersStubReturn: func(brokers []string) {

		},
		WithConsumerIDStubReturn: func(consumerID string) {

		},
		WithAuthStubReturn: func(authMechanism sasl.Mechanism) {

		},
	})
)

func TestKafkaConsumer_Consume(t *testing.T) {
	ctx := context.Background()

	t.Run("test when consumer is successful", func(t *testing.T) {

		consumerStub.WithBrokers([]string{"192.168.1.60:9094"}).
			WithConsumerID("test-consumerID-1").
			WithTopic("test-topic").
			WithAuth(plain.Mechanism{
				Username: "test",
				Password: "test",
			})

		msg, err := consumerStub.Consume(context.Background())
		if err != nil {
			log.Println(err)
		}
		assert.NotNil(t, msg)

	})

	t.Run("test when consumer panics", func(t *testing.T) {

		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()

		msg, err := consumerStubPanic.Consume(ctx)
		fmt.Println(msg, err)

	})

}

func TestKafkaConsumer_ConsumeAsync(t *testing.T) {

	t.Run("test when consumeAsync is success", func(t *testing.T) {
		msgBus := make(chan *kafka.Message)

		err := consumerStub.ConsumeAsync(context.Background(), msgBus, time.Second)

		assert.NoError(t, err)

	})

	t.Run("test when consumeAsync panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()

		msgBus := make(chan *kafka.Message)

		err := consumerStub.ConsumeAsync(context.Background(), msgBus, time.Second)

		assert.NoError(t, err)

	})
}
