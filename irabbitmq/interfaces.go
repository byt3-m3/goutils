package irabbitmq

type Connection interface {
	Close() error

	IsClosed() bool
}
