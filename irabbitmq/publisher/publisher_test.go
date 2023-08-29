package publisher

import (
	"context"
	"github.com/byt3-m3/goutils/env_utils"
	"github.com/byt3-m3/goutils/vars"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
)

func TestNewPublisher(t *testing.T) {
	pub := New().
		WithAMQPUrl(env_utils.GetEnvStrict("AMQP_URL")).
		WithPlainAuth(env_utils.GetEnvStrict("AMQP_USERNAME"), env_utils.GetEnvStrict("AMQP_PASSWORD")).
		WithVHost("/")

	err := pub.Publish(context.Background(), &PublishInput{
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
