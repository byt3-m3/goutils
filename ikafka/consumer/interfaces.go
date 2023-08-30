package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"time"
)

type OptionsSetter interface {
	WithReader(reader *kafka.Reader) Consumer
	WithTopic(topic string) Consumer
	WithBrokers(brokers []string) Consumer
	WithConsumerID(id string) Consumer
	WithAuth(authMechanism sasl.Mechanism) Consumer
}

type Consumer interface {
	OptionsSetter
	ConsumeAsync(ctx context.Context, msgBus chan *kafka.Message, tickerRate time.Duration) error
	Consume(ctx context.Context) (*kafka.Message, error)
}
