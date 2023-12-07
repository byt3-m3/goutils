package irabbitmq

import "github.com/rabbitmq/amqp091-go"

type Connection interface {
	Close() error

	IsClosed() bool
}

type Channel interface {
	Close() error

	Qos(prefetchCount, prefetchSize int, global bool) error

	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error)
}
