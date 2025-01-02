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
	channel    *amqp.Channel
	queue      *amqp.Queue
}

func (r *RabbitMq) Connect() (*amqp.Connection, error) {
	// Server will manage this
	return nil, nil
}

func (r *RabbitMq) Setup() error {
	var err error
	r.channel, err = r.Connection.Channel()
	if err != nil {
		return err
	}
	queue, err := r.channel.QueueDeclare(
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
	r.queue = &queue

	return nil
}

func (r *RabbitMq) Push(ctx context.Context, task common.MqTask) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(ctx,
		"",
		r.queue.Name,
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
