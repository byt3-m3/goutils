package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/env_utils"
	"log"
	"testing"
)

func TestNewConsumer(t *testing.T) {
	c := NewConsumer(
		WithAMQPUrl(env_utils.GetEnvStrict("AMQP_URL")),
		WithPlainAuth(env_utils.GetEnvStrict("AMQP_USERNAME"), env_utils.GetEnvStrict("AMQP_PASSWORD")),
		WithConsumerID("test-consumer-id"),
		WithVhost("/"),
		WithPreFetchCount(1),
	)

	delivery, err := c.Consume(context.Background(), "test-queue")
	if err != nil {
		log.Fatalln(err)
	}

	for msg := range delivery {
		msg.Ack(false)
	}

}
