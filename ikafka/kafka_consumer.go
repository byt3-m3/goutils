package ikafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type kafkaConsumer struct {
	reader *kafka.Reader
}

type KafkaConsumerConfig struct {
	Topic      string
	Brokers    []string
	ConsumerID string
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

func ProvideKafkaConsumer(cfg *KafkaConsumerConfig) IConsumer {
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
