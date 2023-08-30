package publisher

import (
	"context"
	"github.com/byt3-m3/goutils/env_utils"
	"github.com/byt3-m3/goutils/vars"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
)

var (
	livePublisher = New().
			WithAMQPUrl(env_utils.GetEnv("AMQP_URL", "")).
			WithPlainAuth(env_utils.GetEnv("AMQP_USERNAME", ""), env_utils.GetEnv("AMQP_PASSWORD", "")).
			WithVHost("/")

	stubPublisher = NewStubRabbitMQPublisher(&NewStubRabbitMQPublisherInput{
		PublishStubReturn: func(ctx context.Context, input *PublishInput) error {
			return nil
		},
		GetConnectionStubReturn: func() *amqp091.Connection {
			return nil
		},
		ResetConnectionStubReturn: func() error {
			return nil
		},
		IsClosedStubReturn: func() bool {
			return true
		},
		WithAMQPUrlStubReturn: func(url string) {

		},
		WithVHostStubReturn: func(vhost string) {

		},
		WithLoggerStubReturn: func(logger *logrus.Logger) {

		},
		WithNoAuthStubReturn: func() {

		},
		WithPlainAuthStubReturn: func(username, password string) {

		},
	})
)

func TestNewPublisher(t *testing.T) {

	err := stubPublisher.Publish(context.Background(), &PublishInput{
		MessageID:     primitive.NewObjectID().Hex(),
		Exchange:      "test-exchange",
		RoutingKey:    "",
		ContentType:   vars.ContentTypeText,
		CorrelationId: "",
		Headers:       nil,
		Data:          []byte("this is a test message"),
	})

	log.Println(err)

}
