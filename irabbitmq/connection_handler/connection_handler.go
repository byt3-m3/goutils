package connection_handler

import (
	"github.com/byt3-m3/goutils/irabbitmq"
	"github.com/byt3-m3/goutils/logging"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type ConnectionHandler interface {
	GetConnection() irabbitmq.Connection
	MustGetChannel() irabbitmq.Channel
}

type connectionHandler struct {
	amqpUrl        string
	amqpVHost      string
	amqpAuth       amqp091.Authentication
	connection     *amqp091.Connection
	currentChannel *amqp091.Channel
	logger         *log.Logger
}

type NewConnectionHandlerOptions func(h *connectionHandler)

func WithAMQPUrl(url string) NewConnectionHandlerOptions {
	return func(h *connectionHandler) {
		h.amqpUrl = url
	}
}

func WithVHost(vhost string) NewConnectionHandlerOptions {
	return func(h *connectionHandler) {
		h.amqpVHost = vhost
	}
}

func WithPlainAuth(username, password string) NewConnectionHandlerOptions {

	return func(h *connectionHandler) {
		h.amqpAuth = &amqp091.PlainAuth{
			Username: username,
			Password: password,
		}
	}
}

func NewConnectionHandler(opts ...NewConnectionHandlerOptions) ConnectionHandler {
	h := &connectionHandler{}

	for _, opt := range opts {
		opt(h)
	}

	h.MustValidate()
	return h
}

func (c *connectionHandler) MustValidate() {

	if c.amqpUrl == "" {
		panic("AMQPUrl not set")

	}

	if c.logger == nil {
		c.logger = logging.NewLogger()
	}

	if c.connection == nil {
		if err := c.ResetConnection(); err != nil {
			panic(err)
		}
	}

}

func (c *connectionHandler) GetConnection() irabbitmq.Connection {
	if c.connection.IsClosed() {
		if err := c.ResetConnection(); err != nil {
			panic(err)
		}
	}

	return c.connection
}

func (c *connectionHandler) MustGetChannel() irabbitmq.Channel {
	if !c.currentChannel.IsClosed() {
		return c.currentChannel
	}

	if c.connection.IsClosed() {
		if err := c.ResetConnection(); err != nil {
			panic(err)
		}
	}

	ch, err := c.connection.Channel()
	if err != nil {
		panic(err)
	}

	c.currentChannel = ch

	return c.currentChannel

}

func (c *connectionHandler) ResetConnection() error {
	cConn, err := amqp091.DialConfig(c.amqpUrl, amqp091.Config{
		SASL:            []amqp091.Authentication{c.amqpAuth},
		Vhost:           c.amqpVHost,
		ChannelMax:      0,
		FrameSize:       0,
		Heartbeat:       0,
		TLSClientConfig: nil,
		Properties:      nil,
		Locale:          "",
		Dial:            nil,
	})

	if err != nil {
		c.logger.Error(err)
		return err
	}

	c.connection = cConn

	return nil
}
