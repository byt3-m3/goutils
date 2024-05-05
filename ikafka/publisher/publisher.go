package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
)

type publisher struct {
	//dialer        *kafka.Dialer
	writer        *kafka.Writer
	brokerAddress string
	topic         string
}

func (p publisher) PublishMessage(ctx context.Context, msg *kafka.Message) error {

	return p.writer.WriteMessages(ctx, *msg)
}

func NewWithAuth(ctx context.Context, brokerAddr, topic, clientId string, auth sasl.Mechanism) Publisher {

	writer := &kafka.Writer{
		Addr:            kafka.TCP(brokerAddr),
		Topic:           topic,
		Balancer:        &kafka.LeastBytes{},
		MaxAttempts:     0,
		WriteBackoffMin: 0,
		WriteBackoffMax: 0,
		BatchSize:       0,
		BatchBytes:      0,
		BatchTimeout:    0,
		ReadTimeout:     0,
		WriteTimeout:    0,
		RequiredAcks:    0,
		Async:           false,
		Completion:      nil,
		Compression:     0,
		Logger:          nil,
		ErrorLogger:     nil,
		Transport: &kafka.Transport{
			Dial:        nil,
			DialTimeout: 0,
			IdleTimeout: 0,
			MetadataTTL: 0,
			ClientID:    clientId,
			TLS:         nil,
			SASL:        auth,
			Resolver:    nil,
			Context:     ctx,
		},
		AllowAutoTopicCreation: true,
	}

	return publisher{
		//dialer:        dialer,
		brokerAddress: brokerAddr,
		topic:         topic,
		writer:        writer,
	}
}

func NewWithOutAuth(ctx context.Context, brokerAddr, topic, clientId string) Publisher {

	writer := &kafka.Writer{
		Addr:            kafka.TCP(brokerAddr),
		Topic:           topic,
		Balancer:        &kafka.LeastBytes{},
		MaxAttempts:     0,
		WriteBackoffMin: 0,
		WriteBackoffMax: 0,
		BatchSize:       0,
		BatchBytes:      0,
		BatchTimeout:    0,
		ReadTimeout:     0,
		WriteTimeout:    0,
		RequiredAcks:    0,
		Async:           false,
		Completion:      nil,
		Compression:     0,
		Logger:          nil,
		ErrorLogger:     nil,
		Transport: &kafka.Transport{
			Dial:        nil,
			DialTimeout: 0,
			IdleTimeout: 0,
			MetadataTTL: 0,
			ClientID:    clientId,
			TLS:         nil,
			Resolver:    nil,
			Context:     ctx,
		},
		AllowAutoTopicCreation: true,
	}

	return publisher{
		brokerAddress: brokerAddr,
		topic:         topic,
		writer:        writer,
	}
}

//
//func (p publisher) dial() (*kafka.Conn, error) {
//
//	return p.dialer.DialLeader(context.Background(), "tcp", p.brokerAddress, p.topic, p.partition)
//
//}
