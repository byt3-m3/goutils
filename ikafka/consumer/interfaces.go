package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Consumer interface {
	AsyncConsumer
	Consume(ctx context.Context) (*kafka.Message, error)
}

type ConsumeAsyncInput struct {
	MsgChan    chan *kafka.Message
	ErrorChan  chan error
	TickerRate time.Duration
}

type AsyncConsumer interface {
	ConsumeAsync(ctx context.Context, input *ConsumeAsyncInput) error
}

type SyncConsumer interface {
	Consume(ctx context.Context) (*kafka.Message, error)
}
