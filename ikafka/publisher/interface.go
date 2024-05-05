package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type (
	Publisher interface {
		PublishMessage(ctx context.Context, msg *kafka.Message) error
	}
)
