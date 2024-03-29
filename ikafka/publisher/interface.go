package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"log/slog"
)

type (
	OptionSetter interface {
		WithAuth(authMechanism sasl.Mechanism) Publisher

		WithTopic(topic string) Publisher

		WithBroker(broker string) Publisher

		WithKafkaConn(conn *kafka.Conn) Publisher

		WithLogger(logger *slog.Logger) Publisher
	}

	OptionValidator interface {
		MustValidate()
	}

	Publisher interface {
		OptionSetter
		PublishMessage(ctx context.Context, msg *kafka.Message) (int, error)
	}
)
