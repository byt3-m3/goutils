package consumer

import (
	"context"
	"github.com/byt3-m3/goutils/logging"
	"github.com/byt3-m3/goutils/vars"
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
	if len(consumer.brokers) == 0 {
		panic("no brokers set, use WithBrokers")

	}

	if consumer.topic == "" {
		panic("no topic set, use WithTopic")

	}

	if consumer.consumerID == "" {
		panic("no ConsumerID set, use WithConsumerID")

	}

	if consumer.logger == nil {
		consumer.logger = logging.NewLogger()
	}
	if consumer.authMechanism == nil {
		consumer.authMechanism = plain.Mechanism{
			Username: vars.KafkaUsername,
			Password: vars.KafkaPassword,
		}
	}

	if consumer.reader == nil {
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
	c.logger.Info("starting async consumer")
	for {
		select {
		case <-ticker.C:

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
		c.logger.Error("error reading message")
		return nil, err
	}
	return &msg, nil
}
