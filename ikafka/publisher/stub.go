package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"log/slog"
)

type StubKafkaPublisher struct {
	PublishMessageStubReturn func(ctx context.Context, msg *kafka.Message) (int, error)
	WithAuthStubReturn       func(authMechanism sasl.Mechanism)
	WithTopicStubReturn      func(topic string)
	WithBrokerStubReturn     func(broker string)
	WithKafkaConnStubReturn  func(conn *kafka.Conn)
	WithLoggerStubReturn     func(logger *slog.Logger)
}

type NewStubKafkaPublisherInput struct {
	PublishMessageStubReturn func(ctx context.Context, msg *kafka.Message) (int, error)
	WithAuthStubReturn       func(authMechanism sasl.Mechanism)
	WithTopicStubReturn      func(topic string)
	WithBrokerStubReturn     func(broker string)
	WithKafkaConnStubReturn  func(conn *kafka.Conn)
	WithLoggerStubReturn     func(logger *slog.Logger)
}

func NewStubKafkaPublisher(input *NewStubKafkaPublisherInput) Publisher {

	return &StubKafkaPublisher{
		PublishMessageStubReturn: input.PublishMessageStubReturn,
		WithAuthStubReturn:       input.WithAuthStubReturn,
		WithTopicStubReturn:      input.WithTopicStubReturn,
		WithBrokerStubReturn:     input.WithBrokerStubReturn,
		WithKafkaConnStubReturn:  input.WithKafkaConnStubReturn,
		WithLoggerStubReturn:     input.WithLoggerStubReturn,
	}
}

func (s *StubKafkaPublisher) WithLogger(logger *slog.Logger) Publisher {
	s.WithLoggerStubReturn(logger)
	return s
}

func (s *StubKafkaPublisher) WithAuth(authMechanism sasl.Mechanism) Publisher {
	s.WithAuthStubReturn(authMechanism)
	return s
}

func (s *StubKafkaPublisher) WithTopic(topic string) Publisher {
	s.WithTopicStubReturn(topic)
	return s
}

func (s *StubKafkaPublisher) WithBroker(broker string) Publisher {
	s.WithBrokerStubReturn(broker)
	return s

}

func (s *StubKafkaPublisher) WithKafkaConn(conn *kafka.Conn) Publisher {
	s.WithKafkaConnStubReturn(conn)
	return s
}

type PublishMessageStubReturn struct {
	Count int
	Error error
}

func (s *StubKafkaPublisher) PublishMessage(ctx context.Context, msg *kafka.Message) (int, error) {
	return s.PublishMessageStubReturn(ctx, msg)
}
