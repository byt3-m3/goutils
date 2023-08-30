package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/env_utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	log "github.com/sirupsen/logrus"

	"time"
)

type kafkaConsumer struct {
	reader        *kafka.Reader
	topic         string
	brokers       []string
	consumerID    string
	authMechanism sasl.Mechanism
	logger        *log.Logger
}

func New() Consumer {
	return &kafkaConsumer{}
}

func (c *kafkaConsumer) WithLogger(logger *log.Logger) Consumer {
	c.logger = logger
	return c
}

func (c *kafkaConsumer) WithReader(reader *kafka.Reader) Consumer {
	c.reader = reader
	return c
}

func (c *kafkaConsumer) WithTopic(topic string) Consumer {
	c.topic = topic
	return c
}

func (c *kafkaConsumer) WithBrokers(brokers []string) Consumer {
	c.brokers = brokers
	return c
}

func (c *kafkaConsumer) WithConsumerID(id string) Consumer {
	c.consumerID = id
	return c
}

func (c *kafkaConsumer) WithAuth(authMechanism sasl.Mechanism) Consumer {
	c.authMechanism = authMechanism
	return c
}

func MustValidate(consumer *kafkaConsumer) bool {

	switch {
	case len(consumer.brokers) == 0:
		panic("no brokers set, use WithBrokers")

	case consumer.topic == "":
		panic("no topic set, use WithTopic")

	case consumer.consumerID == "":
		panic("no ConsumerID set, use WithConsumerID")

	case consumer.authMechanism == nil:
		consumer.authMechanism = plain.Mechanism{
			Username: env_utils.GetEnvStrict("KAFKA_USERNAME"),
			Password: env_utils.GetEnvStrict("KAFKA_PASSWORD"),
		}

	case consumer.reader == nil:
		dialer := &kafka.Dialer{

			SASLMechanism: consumer.authMechanism,
		}

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        consumer.brokers,
			GroupID:        consumer.consumerID,
			Topic:          consumer.topic,
			Partition:      0,
			QueueCapacity:  0,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			Dialer:         dialer,
			IsolationLevel: 0,
			MaxAttempts:    3,
		})

		consumer.reader = reader

	}

	return true

}

type ConsumeAsyncInput struct {
	MsgChan    chan *kafka.Message
	ErrorChan  chan error
	TickerRate time.Duration
}

func (c *kafkaConsumer) ConsumeAsync(ctx context.Context, input *ConsumeAsyncInput) error {
	MustValidate(c)

	ticker := time.NewTicker(input.TickerRate)

	for {
		select {
		case <-ticker.C:

			log.Println("executing read")
			msg, err := c.reader.ReadMessage(ctx)

			if err != nil {
				c.logger.Error("error reading message", err)
				input.ErrorChan <- err
			}

			input.MsgChan <- &msg
		}

	}
}

func (c *kafkaConsumer) Consume(ctx context.Context) (*kafka.Message, error) {
	MustValidate(c)
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		log.Fatalln("error reading message", err)
	}
	return &msg, nil
}
