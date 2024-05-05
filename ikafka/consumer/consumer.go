package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"log/slog"
)

type consumer struct {
	reader        *kafka.Reader
	topic         string
	brokers       []string
	consumerID    string
	authMechanism sasl.Mechanism
	logger        *slog.Logger
}

func (c consumer) Consume(ctx context.Context) (*kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		c.logger.Error("error reading message",
			slog.Any("error", err),
		)
		return nil, err
	}
	return &msg, nil
}

func NewWithAuth(brokers []string, topic, consumerID string, auth sasl.Mechanism, logger *slog.Logger) SyncConsumer {

	dialer := &kafka.Dialer{

		SASLMechanism: auth,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        consumerID,
		Topic:          topic,
		Partition:      0,
		QueueCapacity:  0,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		Dialer:         dialer,
		IsolationLevel: 0,
		MaxAttempts:    3,
	})

	return consumer{
		reader:        reader,
		topic:         topic,
		brokers:       brokers,
		consumerID:    consumerID,
		authMechanism: auth,
		logger:        logger,
	}
}

func NewWithOutAuth(brokers []string, topic, consumerID string, logger *slog.Logger) SyncConsumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        consumerID,
		Topic:          topic,
		Partition:      0,
		QueueCapacity:  0,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		Dialer:         nil,
		IsolationLevel: 0,
		MaxAttempts:    3,
	})

	return consumer{
		reader:        reader,
		topic:         topic,
		brokers:       brokers,
		consumerID:    consumerID,
		authMechanism: nil,
		logger:        logger,
	}
}
