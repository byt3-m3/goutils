package admin

type AMQPExchangeType string

var (
	AMQPExchangeTypeDirect  AMQPExchangeType = "direct"
	AMQPExchangeTypeFanout  AMQPExchangeType = "fanout"
	AMQPExchangeTypeHeaders AMQPExchangeType = "headers"
	AMQPExchangeTypeTopic   AMQPExchangeType = "topic"
)
