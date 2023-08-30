package publisher

import (
	"context"
	"github.com/byt3-m3/goutils/env_utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	log "github.com/sirupsen/logrus"

	"time"
)

type kafkaPublisher struct {
	topic         string
	brokerAddr    string
	conn          *kafka.Conn
	partition     int
	authMechanism sasl.Mechanism
	logger        *log.Logger
}

func NewPublisher() Publisher {
	return &kafkaPublisher{}
}

func (p *kafkaPublisher) WithLogger(logger *log.Logger) Publisher {
	p.logger = logger
	return p
}

func (p *kafkaPublisher) WithAuth(authMechanism sasl.Mechanism) Publisher {
	p.authMechanism = authMechanism
	return p
}

func (p *kafkaPublisher) WithTopic(topic string) Publisher {
	p.topic = topic
	return p

}

func (p *kafkaPublisher) WithBroker(broker string) Publisher {
	p.brokerAddr = broker
	return p
}

func (p *kafkaPublisher) WithKafkaConn(conn *kafka.Conn) Publisher {
	p.conn = conn
	return p
}

func (p *kafkaPublisher) MustValidate() {
	switch {
	case p.brokerAddr == "":
		panic("broker address not set, use WithBroker")

	case p.topic == "":
		panic("topic not set, user WithTopic")

	case p.conn == nil:

		switch p.authMechanism.(type) {
		case nil:
			p.authMechanism = plain.Mechanism{
				Username: env_utils.GetEnvStrict("KAFKA_USERNAME"),
				Password: env_utils.GetEnvStrict("KAFKA_PASSWORD"),
			}

		}

		dialer := kafka.Dialer{
			ClientID:        "",
			DialFunc:        nil,
			Timeout:         0,
			Deadline:        time.Time{},
			LocalAddr:       nil,
			DualStack:       false,
			FallbackDelay:   0,
			KeepAlive:       0,
			Resolver:        nil,
			TLS:             nil,
			SASLMechanism:   p.authMechanism,
			TransactionalID: "",
		}
		conn, err := dialer.DialLeader(context.Background(), "tcp", p.brokerAddr, p.topic, p.partition)
		if err != nil {
			panic(err)
		}

		p.conn = conn
	}
}

func (p *kafkaPublisher) PublishMessage(ctx context.Context, msg *kafka.Message) (int, error) {
	written, err := p.conn.WriteMessages(*msg)
	if err != nil {
		p.logger.Error(err)
		return 0, err
	}
	return written, nil
}
