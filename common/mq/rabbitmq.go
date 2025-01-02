package mq

import (
	common "common"
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QUEUE_NAME = "transcode_queue"

type RabbitMq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}

func (r *RabbitMq) Connect() (*amqp.Connection, error) {
	// Server will manage this
	return nil, nil
}

func (r *RabbitMq) Setup() error {
	var err error
	r.Channel, err = r.Connection.Channel()
	if err != nil {
		return err
	}
	queue, err := r.Channel.QueueDeclare(
		QUEUE_NAME, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}
	r.Queue = &queue

	return nil
}

func (r *RabbitMq) Push(ctx context.Context, task common.MqTask) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = r.Channel.PublishWithContext(ctx,
		"",
		r.Queue.Name,
		true,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
