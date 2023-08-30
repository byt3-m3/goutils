package consumer

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var (
	stubConsumer = NewStubRabbitMQConsumer(&NewStubRabbitMQConsumerInput{
		ConsumeReturn: func(ctx context.Context, queue string) (<-chan amqp091.Delivery, error) {
			data := make(chan amqp091.Delivery, 100)

			return data, nil
		},
		GetConnectionReturn: func() *amqp091.Connection {

			return nil
		},
		IsClosedReturn: func() bool {
			return true
		},
		GetActiveChannelStubReturn: func() *amqp091.Channel {
			return nil
		},
		WithAMQPUrlStubReturn: func(url string) {

		},
		WithConsumerIDStubReturn: func(id string) {

		},
		WithVHostStubReturn: func(vhost string) {

		},
		WithPlainAuthStubReturn: func(username, password string) {

		},
		WithPreFetchCountStubReturn: func(count int) {

		},
		WithLoggerStubReturn: func(logger *logrus.Logger) {

		},
		ResetConnectionStubReturn: func() {

		},
	})
)

func TestNewConsumer(t *testing.T) {

	delivery, err := stubConsumer.Consume(context.Background(), "test-queue")
	if err != nil {
		log.Fatalln(err)
	}
	assert.NotNil(t, delivery)
}
