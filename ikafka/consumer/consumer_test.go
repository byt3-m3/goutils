package consumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestConsumer_Consume(t *testing.T) {
	consumer := NewWithAuth([]string{"kafka1.baxterhome.net:9094"}, "test-topic", "tester", plain.Mechanism{
		Username: "kafka",
		Password: "kafka",
	}, slog.Default())

	msg, err := consumer.Consume(context.Background())
	fmt.Println(msg)
	assert.NoError(t, err)
	assert.NotNil(t, msg)
}
