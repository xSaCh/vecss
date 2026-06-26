package mq

import (
	"context"
	"encoding/json"
	"maps"
	"vecss/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QUEUE_NAME = "transcode_queue"

type rabbitMessage struct {
	delivery amqp.Delivery
}

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

func (r *RabbitMq) Push(ctx context.Context, task domain.MqTask) error {
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

func (r *RabbitMq) Consume(ctx context.Context) (<-chan MqMessage, error) {
	deliveries, err := r.Channel.ConsumeWithContext(ctx, r.Queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// Bridge the amqp.Delivery channel into a domain.Message channel
	out := make(chan MqMessage)
	go func() {
		defer close(out)
		for d := range deliveries {
			select {
			case out <- &rabbitMessage{delivery: d}:
			case <-ctx.Done():
				return
			}

		}
	}()
	return out, nil
}

// Implementing MqMessage interface for rabbitMessage

func (m *rabbitMessage) Body() []byte {
	return m.delivery.Body
}

func (m *rabbitMessage) Ack() error {
	return m.delivery.Ack(false)
}

func (m *rabbitMessage) Nack(requeue bool) error {
	return m.delivery.Nack(false, requeue)
}

func (m *rabbitMessage) Headers() map[string]any {
	headers := make(map[string]any)
	maps.Copy(headers, m.delivery.Headers)
	return headers
}
