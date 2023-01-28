package ikafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type (
	IPublisher interface {
		PublishMessage(ctx context.Context, msg *kafka.Message) (int, error)
	}

	publisher struct {
		Topic      string
		BrokerAddr string
		Conn       *kafka.Conn
	}

	KafkaPublisherConfig struct {
		Topic      string
		BrokerAddr string
		Partition  int
	}
)

func ProvideKafkaPublisher(cfg *KafkaPublisherConfig) IPublisher {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.BrokerAddr, cfg.Topic, cfg.Partition)
	if err != nil {
		log.Fatalln(err)
	}
	return &publisher{
		Topic:      cfg.Topic,
		BrokerAddr: cfg.BrokerAddr,
		Conn:       conn,
	}

}

func (p *publisher) PublishMessage(ctx context.Context, msg *kafka.Message) (int, error) {
	written, err := p.Conn.WriteMessages(*msg)
	if err != nil {
		log.Println("error writing message")
		return 0, err
	}
	return written, nil
}
