package publisher

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"

	"time"
)

type publisher struct {
	amqpUrl  string
	vHost    string
	logger   *log.Logger
	conn     *amqp091.Connection
	amqpAuth amqp091.Authentication
}

func New() RabbitMQPublisher {
	p := &publisher{}

	return p
}

func (p *publisher) MustValidate() {
	switch {
	case p.amqpUrl == "":
		panic("amqpURL not set")

	case p.conn == nil:
		if err := p.ResetConnection(); err != nil {
			p.logger.Fatalln(err)
		}

	case p.logger == nil:
		p.logger = log.New()
	}
}

func (p *publisher) WithAMQPUrl(url string) RabbitMQPublisher {
	p.amqpUrl = url
	return p
}

func (p *publisher) WithVHost(vhost string) RabbitMQPublisher {
	p.vHost = vhost
	return p
}

func (p *publisher) WithLogger(log *log.Logger) RabbitMQPublisher {
	p.logger = log
	return p
}

func (p *publisher) WithNoAuth() RabbitMQPublisher {
	p.amqpAuth = nil
	return p
}

func (p *publisher) WithPlainAuth(username, password string) RabbitMQPublisher {
	p.amqpAuth = &amqp091.PlainAuth{
		Username: username,
		Password: password,
	}
	return p
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
