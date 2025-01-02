package mq

import (
	common "common"
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter interface {
	Setup() error
	Push(ctx context.Context, task common.MqTask) error
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
