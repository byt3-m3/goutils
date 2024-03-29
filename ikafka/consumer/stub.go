package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	log "github.com/sirupsen/logrus"
)

type StubKafkaConsumer struct {
	ConsumeAsyncStubReturn   func(ctx context.Context, input *ConsumeAsyncInput) error
	ConsumeStubReturn        func(ctx context.Context) (*kafka.Message, error)
	WithReaderStubReturn     func(reader *kafka.Reader)
	WithTopicStubReturn      func(topic string)
	WithBrokersStubReturn    func(brokers []string)
	WithConsumerIDStubReturn func(consumerID string)
	WithAuthStubReturn       func(authMechanism sasl.Mechanism)
	WithLoggerStubReturn     func(logger *log.Logger)
}

func (s *StubKafkaConsumer) WithLogger(logger *log.Logger) Consumer {
	s.WithLoggerStubReturn(logger)
	return s
}

func (s *StubKafkaConsumer) WithReader(reader *kafka.Reader) Consumer {
	s.WithReaderStubReturn(reader)
	return s
}

func (s *StubKafkaConsumer) WithTopic(topic string) Consumer {
	s.WithTopicStubReturn(topic)
	return s
}

func (s *StubKafkaConsumer) WithBrokers(brokers []string) Consumer {
	s.WithBrokersStubReturn(brokers)
	return s
}

func (s *StubKafkaConsumer) WithConsumerID(id string) Consumer {
	s.WithConsumerIDStubReturn(id)
	return s
}

func (s *StubKafkaConsumer) WithAuth(authMechanism sasl.Mechanism) Consumer {
	s.WithAuthStubReturn(authMechanism)
	return s
}

type NewStubInput struct {
	ConsumeAsyncStubReturn   func(ctx context.Context, input *ConsumeAsyncInput) error
	ConsumeStubReturn        func(ctx context.Context) (*kafka.Message, error)
	WithReaderStubReturn     func(reader *kafka.Reader)
	WithTopicStubReturn      func(topic string)
	WithBrokersStubReturn    func(brokers []string)
	WithConsumerIDStubReturn func(consumerID string)
	WithAuthStubReturn       func(authMechanism sasl.Mechanism)
}

func NewStub(input NewStubInput) Consumer {

	return &StubKafkaConsumer{
		ConsumeAsyncStubReturn:   input.ConsumeAsyncStubReturn,
		ConsumeStubReturn:        input.ConsumeStubReturn,
		WithReaderStubReturn:     input.WithReaderStubReturn,
		WithTopicStubReturn:      input.WithTopicStubReturn,
		WithBrokersStubReturn:    input.WithBrokersStubReturn,
		WithConsumerIDStubReturn: input.WithConsumerIDStubReturn,
		WithAuthStubReturn:       input.WithAuthStubReturn,
	}
}

type ConsumeAsyncStubReturn struct {
	Error error
}

func (s *StubKafkaConsumer) ConsumeAsync(ctx context.Context, input *ConsumeAsyncInput) error {
	return s.ConsumeAsyncStubReturn(ctx, input)
}

type ConsumeStubReturn struct {
	Message *kafka.Message
	Error   error
}

func (s *StubKafkaConsumer) Consume(ctx context.Context) (*kafka.Message, error) {
	return s.ConsumeStubReturn(ctx)
}
