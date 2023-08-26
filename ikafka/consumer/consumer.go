package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Consumer interface {
	ConsumeAsync(ctx context.Context, msgBus chan *kafka.Message, tickerRate time.Duration) error
	Consume(ctx context.Context) (*kafka.Message, error)
}
