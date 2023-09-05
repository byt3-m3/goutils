package admin_client

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"testing"
)

var (
	stubAdminClient = NewStubClient(
		&NewStubClientInput{
			CreateQueueStubReturn: func(ctx context.Context, input *CreateQueueInput) (*amqp091.Queue, error) {
				return &amqp091.Queue{}, nil
			},
			CreateExchangeStubReturn: func(ctx context.Context, input *CreateExchangeInput) error {
				return nil
			},
			BindQueueStubReturn: func(ctx context.Context, input *BindQueueInput) error {
				return nil
			},
			DeleteQueueStubReturn: func(ctx context.Context, input *DeleteQueueInput) error {
				return nil
			},
			DeleteExchangeStubReturn: func(ctx context.Context, input *DeleteExchangeInput) error {
				return nil
			},
			GetConnectionStubReturn: func() *amqp091.Connection {

				return &amqp091.Connection{}
			},
			WithAMQPUrlStubReturn: func(url string) {

			},
			WithVHostStubReturn: func(vhost string) {

			},
			WithPlainAuthStubReturn: func(username, password string) {

			},
			WithConnectionStubReturn: func(conn *amqp091.Connection) {

			},
			ValidateClientStubReturn: func(client *adminClient) bool {
				return true
			},
		},
	)

	ac = New().WithVHost("/").WithPlainAuth("user", "pass").WithAMQPUrl("amqp://localhost")
)

func TestNewAdminClient(t *testing.T) {

	q, err := stubAdminClient.CreateQueue(context.Background(), &CreateQueueInput{
		QueName:       "test-queue",
		IsDurable:     false,
		IsExclusive:   false,
		CanAutoDelete: false,
		CanNoWait:     false,
	})

	log.Print(q, err)

	err = stubAdminClient.CreateExchange(context.Background(), &CreateExchangeInput{
		ExchangeName:  "test-exchange",
		ExchangeType:  AMQPExchangeTypeDirect,
		IsDurable:     false,
		IsInternal:    false,
		CanAutoDelete: false,
		CanNoWait:     false,
	})

	log.Print(err)

	err = stubAdminClient.BindQueue(nil, &BindQueueInput{
		QueName:   "test-queue",
		Key:       "",
		Exchange:  "test-exchange",
		CanNoWait: false,
		Args:      nil,
	})

	log.Print(err)

}

func TestAdminClient(t *testing.T) {

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
