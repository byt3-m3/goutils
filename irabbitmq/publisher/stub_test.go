package publisher

import "testing"

func TestStubRabbitMQPublisher(t *testing.T) {
	stub := StubRabbitMQPublisher{
		PublishStubReturn:         nil,
		GetConnectionStubReturn:   nil,
		ResetConnectionStubReturn: nil,
		IsClosedStubReturn:        nil,
		WithAMQPUrlStubReturn: func(url string) {

		},
		WithVHostStubReturn:  nil,
		WithLoggerStubReturn: nil,
		WithNoAuthStubReturn: nil,
		WithPlainAuthStubReturn: func(username, password string) {

		},
		MustValidateStubReturn: func() {

		},
	}

	stub.WithPlainAuth("test", "password").WithAMQPUrl("amqp://url.com")
}
