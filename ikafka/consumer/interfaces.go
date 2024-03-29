package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	log "github.com/sirupsen/logrus"
)

type OptionsSetter interface {
	WithReader(reader *kafka.Reader) Consumer
	WithTopic(topic string) Consumer
	WithBrokers(brokers []string) Consumer
	WithConsumerID(id string) Consumer
	WithAuth(authMechanism sasl.Mechanism) Consumer

	WithLogger(logger *log.Logger) Consumer
}

type Consumer interface {
	OptionsSetter
	ConsumeAsync(ctx context.Context, input *ConsumeAsyncInput) error
	Consume(ctx context.Context) (*kafka.Message, error)
}
