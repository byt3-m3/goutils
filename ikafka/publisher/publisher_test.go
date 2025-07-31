package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func _TestPublisher_PublishMessage(t *testing.T) {
	p := NewWithAuth(context.Background(), "kafka1.baxterhome.net:9094", "test-topic", "test-client", plain.Mechanism{
		Username: "kafka",
		Password: "kafka",
	})

	err := p.PublishMessage(context.Background(), &kafka.Message{
		Topic:         "new-test-topic",
		Partition:     0,
		Offset:        0,
		HighWaterMark: 0,
		Key:           nil,
		Value:         []byte("test"),
		Headers:       nil,
		Time:          time.Time{},
	})

	assert.NoError(t, err)

}
