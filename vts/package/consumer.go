package vts

import (
	"common"
	"common/mq"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Consumer struct {
	Rbmq *mq.RabbitMq
}

func NewConsumer(rbmq *mq.RabbitMq) *Consumer {
	c := Consumer{
		Rbmq: rbmq,
	}
	err := c.Rbmq.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil

	}
	return &c
}

func (c *Consumer) Listen(ctx context.Context) error {
	tasks, err := c.Rbmq.Channel.ConsumeWithContext(ctx, c.Rbmq.Queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for task := range tasks {
			var MqTask common.MqTask
			json.Unmarshal(task.Body, &MqTask)

			handleTask(MqTask)
		}
	}()
	<-forever
	return nil
}

func handleTask(task common.MqTask) {
	fmt.Printf("Received task: %v\n", task)
	time.Sleep(10 * time.Second)
	fmt.Printf("Done work.....")
}
