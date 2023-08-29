package admin_client

import (
	"context"
	"log"
	"testing"
)

func TestNewAdminClient(t *testing.T) {

	ac := New().
		WithVHost("/").
		WithAMQPUrl("amqp://192.168.1.61").
		WithPlainAuth("cbaxter", "pimpin12")

	q, err := ac.CreateQueue(context.Background(), &CreateQueueInput{
		QueName:       "test-queue",
		IsDurable:     false,
		IsExclusive:   false,
		CanAutoDelete: false,
		CanNoWait:     false,
	})

	log.Print(q, err)

	err = ac.CreateExchange(context.Background(), &CreateExchangeInput{
		ExchangeName:  "test-exchange",
		ExchangeType:  AMQPExchangeTypeDirect,
		IsDurable:     false,
		IsInternal:    false,
		CanAutoDelete: false,
		CanNoWait:     false,
	})

	log.Print(err)

	err = ac.BindQueue(nil, &BindQueueInput{
		QueName:   "test-queue",
		Key:       "",
		Exchange:  "test-exchange",
		CanNoWait: false,
		Args:      nil,
	})

	log.Print(err)

}
