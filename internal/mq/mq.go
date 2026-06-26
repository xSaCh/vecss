package mq

import (
	"context"
	"fmt"
	"vecss/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Abstracts over the broker-specific delivery type
type MqMessage interface {
	Body() []byte
	Ack() error
	Nack(requeue bool) error
	Headers() map[string]interface{}
}

type Emitter interface {
	Setup() error
	Push(ctx context.Context, task domain.MqTask) error
	Consume(ctx context.Context) (<-chan MqMessage, error)
}

func NewRabbitMqEmitter(username, password, url string) (*RabbitMq, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/", username, password, url))

	if err != nil {
		return nil, err
	}
	return &RabbitMq{
		Connection: conn,
	}, nil
}
