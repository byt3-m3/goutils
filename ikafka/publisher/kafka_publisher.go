package publisher

import (
	"context"
	"github.com/byt3-m3/goutils/env_utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"log"
	"time"
)

type kafkaPublisher struct {
	topic      string
	brokerAddr string
	conn       *kafka.Conn
}

type (
	opts struct {
		topic         string
		brokerAddr    string
		partition     int
		conn          *kafka.Conn
		authMechanism sasl.Mechanism
	}

	opt func(opt *opts)
)

var (
	WithTopic = func(topic string) opt {
		return func(opt *opts) {
			opt.topic = topic
		}
	}

	WithBroker = func(address string) opt {
		return func(opt *opts) {
			opt.brokerAddr = address
		}
	}

	WithKafkaConn = func(conn *kafka.Conn) opt {
		return func(opt *opts) {
			opt.conn = conn
		}
	}

	WithPlainAuth = func(authMechanism plain.Mechanism) opt {
		return func(opt *opts) {
			opt.authMechanism = authMechanism
		}
	}
)

func NewPublisher(inputOpts ...opt) (Publisher, error) {
	pOpts := &opts{}

	for _, opt := range inputOpts {
		opt(pOpts)
	}

	if pOpts.conn == nil {
		if pOpts.authMechanism == nil {
			pOpts.authMechanism = plain.Mechanism{
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
			SASLMechanism:   pOpts.authMechanism,
			TransactionalID: "",
		}
		conn, err := dialer.DialLeader(context.Background(), "tcp", pOpts.brokerAddr, pOpts.topic, pOpts.partition)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return &kafkaPublisher{
			topic:      pOpts.topic,
			brokerAddr: pOpts.brokerAddr,
			conn:       conn,
		}, nil
	}

	return &kafkaPublisher{
		topic:      pOpts.topic,
		brokerAddr: pOpts.brokerAddr,
		conn:       pOpts.conn,
	}, nil
}

func (p *kafkaPublisher) PublishMessage(ctx context.Context, msg *kafka.Message) (int, error) {
	written, err := p.conn.WriteMessages(*msg)
	if err != nil {
		log.Println("error writing message")
		return 0, err
	}
	return written, nil
}
