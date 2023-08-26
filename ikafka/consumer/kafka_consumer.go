package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/env_utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"log"
	"time"
)

var (
	WithTopic = func(topic string) opt {

		return func(cfg *options) {
			cfg.Topic = topic
		}
	}

	WithBrokers = func(brokers []string) opt {

		return func(cfg *options) {
			cfg.Brokers = brokers
		}
	}

	WithConsumerID = func(consumerID string) opt {

		return func(cfg *options) {
			cfg.ConsumerID = consumerID
		}
	}

	WithReader = func(reader *kafka.Reader) opt {

		return func(cfg *options) {
			cfg.reader = reader
		}
	}

	WithPlainAuth = func(mechanism plain.Mechanism) opt {

		return func(cfg *options) {
			cfg.authMechanism = mechanism
		}
	}

	WithDefaultReader = func() opt {

		return func(cfg *options) {

			cfg.reader = kafka.NewReader(kafka.ReaderConfig{
				Brokers:       cfg.Brokers,
				GroupID:       cfg.ConsumerID,
				Topic:         cfg.Topic,
				Partition:     0,
				QueueCapacity: 0,
				MinBytes:      10e3,
				MaxBytes:      10e6,

				IsolationLevel: 0,
				MaxAttempts:    3,
			})
		}
	}
)

type opt func(cfg *options)

type kafkaConsumer struct {
	reader *kafka.Reader
}

func NewConsumer(inputOpts ...opt) Consumer {
	consumerOptions := &options{}

	for _, opt := range inputOpts {
		opt(consumerOptions)
	}

	if !validateOptions(consumerOptions) {
		log.Fatalln("failed option validation")
	}

	if consumerOptions.authMechanism == nil {
		consumerOptions.authMechanism = plain.Mechanism{
			Username: env_utils.GetEnvStrict("KAFKA_USERNAME"),
			Password: env_utils.GetEnvStrict("KAFKA_PASSWORD"),
		}
	}

	if consumerOptions.reader == nil {

		dialer := &kafka.Dialer{

			SASLMechanism: consumerOptions.authMechanism,
		}

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        consumerOptions.Brokers,
			GroupID:        consumerOptions.ConsumerID,
			Topic:          consumerOptions.Topic,
			Partition:      0,
			QueueCapacity:  0,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			Dialer:         dialer,
			IsolationLevel: 0,
			MaxAttempts:    3,
		})

		return kafkaConsumer{reader: reader}
	}

	return kafkaConsumer{reader: consumerOptions.reader}
}

func validateOptions(ops *options) bool {
	if len(ops.Brokers) == 0 {
		log.Println("no brokers set, use WithBrokers")
		return false
	}

	if ops.Topic == "" {
		log.Println("no topic set, use WithTopic")

		return false
	}
	if ops.ConsumerID == "" {
		log.Println("no ConsumerID set, use WithConsumerID")
		return false
	}

	return true

}

type options struct {
	Topic         string
	Brokers       []string
	ConsumerID    string
	reader        *kafka.Reader
	authMechanism sasl.Mechanism
}

func (c kafkaConsumer) ConsumeAsync(ctx context.Context, msgBus chan *kafka.Message, tickerRate time.Duration) error {
	ticker := time.NewTicker(tickerRate)

	for {
		select {
		case <-ticker.C:

			log.Println("executing read")
			msg, err := c.reader.ReadMessage(ctx)

			if err != nil {
				log.Fatalln("error reading message", err)
			}

			msgBus <- &msg
		}

	}
}

func (c kafkaConsumer) Consume(ctx context.Context) (*kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		log.Fatalln("error reading message", err)
	}
	return &msg, nil
}

func ProvideKafkaConsumer(cfg *options) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:       cfg.Brokers,
		GroupID:       cfg.ConsumerID,
		Topic:         cfg.Topic,
		Partition:     0,
		QueueCapacity: 0,
		MinBytes:      10e3,
		MaxBytes:      10e6,

		IsolationLevel: 0,
		MaxAttempts:    3,
	})

	return &kafkaConsumer{reader: reader}

}
