package publisher

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type publisherOpt func(pub *publisher)

var (
	WithAMQPUrl = func(url string) publisherOpt {
		return func(pub *publisher) {
			pub.amqpUrl = url
		}
	}

	WithVHost = func(vhost string) publisherOpt {
		return func(pub *publisher) {
			pub.vHost = vhost
		}
	}

	WithPlainAuth = func(username, password string) publisherOpt {
		return func(pub *publisher) {
			pub.amqpAuth = &amqp091.PlainAuth{
				Username: username,
				Password: password,
			}
		}
	}

	WithLogger = func(logger *log.Logger) publisherOpt {
		return func(pub *publisher) {
			pub.logger = logger
		}
	}

	WithNoAuth = func() publisherOpt {
		return func(pub *publisher) {
			pub.amqpAuth = nil
		}
	}
)

type Publisher interface {
	Publish(ctx context.Context, input *PublishInput) error
	GetConnection() *amqp091.Connection
	ResetConnection() error
	IsClosed() bool
}

type publisher struct {
	amqpUrl  string
	vHost    string
	logger   *log.Logger
	conn     *amqp091.Connection
	amqpAuth amqp091.Authentication
}

func NewPublisher(opts ...publisherOpt) Publisher {
	p := &publisher{}

	for _, opt := range opts {
		opt(p)
	}

	if p.logger == nil {
		p.logger = log.Default()
	}

	if !validatePublisher(p) {
		p.logger.Fatalln("failed publisher validation")
	}

	if p.conn == nil {

		if err := p.ResetConnection(); err != nil {
			p.logger.Fatalln(err)
		}
	}

	return p
}

func validatePublisher(p *publisher) bool {
	if p.amqpUrl == "" {
		log.Println("amqpUrl not set, use WithAMQPUrl")

		return false
	}

	return true
}

type PublishInput struct {
	MessageID     string
	Exchange      string
	RoutingKey    string
	ContentType   string
	CorrelationId string
	Headers       amqp091.Table
	Data          []byte
}

func (p *publisher) Publish(ctx context.Context, input *PublishInput) error {
	if p.IsClosed() {
		if err := p.ResetConnection(); err != nil {
			p.logger.Println(err)
			return err
		}
	}

	ch, err := p.conn.Channel()
	if err != nil {
		p.logger.Println(err)
		return err
	}
	defer ch.Close()

	if err := ch.PublishWithContext(ctx, input.Exchange, input.RoutingKey, false, false, amqp091.Publishing{
		Headers:         input.Headers,
		ContentType:     input.ContentType,
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   input.CorrelationId,
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       input.MessageID,
		Timestamp:       time.Time{},

		UserId: "",
		AppId:  "",
		Body:   input.Data,
	}); err != nil {
		p.logger.Println(err)
		return err
	}

	return nil

}

func (p *publisher) GetConnection() *amqp091.Connection {
	return p.conn
}

func (p *publisher) ResetConnection() error {
	conn, err := amqp091.DialConfig(p.amqpUrl, amqp091.Config{
		SASL:            []amqp091.Authentication{p.amqpAuth},
		Vhost:           p.vHost,
		ChannelMax:      0,
		FrameSize:       0,
		Heartbeat:       0,
		TLSClientConfig: nil,
		Properties:      nil,
		Locale:          "",
		Dial:            nil,
	})

	if err != nil {
		p.logger.Println(err)
		return err
	}

	p.conn = conn

	return nil
}

func (p *publisher) IsClosed() bool {
	return p.conn.IsClosed()
}
