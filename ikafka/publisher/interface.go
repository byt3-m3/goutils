package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	log "github.com/sirupsen/logrus"
)

type (
	OptionSetter interface {
		WithAuth(authMechanism sasl.Mechanism) Publisher

		WithTopic(topic string) Publisher

		WithBroker(broker string) Publisher

		WithKafkaConn(conn *kafka.Conn) Publisher

		WithLogger(logger *log.Logger) Publisher
	}

	OptionValidator interface {
		MustValidate()
	}

	Publisher interface {
		OptionSetter
		PublishMessage(ctx context.Context, msg *kafka.Message) (int, error)
	}
)
