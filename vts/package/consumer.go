package vts

import (
	"common"
	"common/mq"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Consumer struct {
	Rbmq       *mq.RabbitMq
	Transcoder Transcoder
}

func NewConsumer(rbmq *mq.RabbitMq, transcoder Transcoder) *Consumer {
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
	c.Transcoder = transcoder
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
			var mqTask common.MqTask
			json.Unmarshal(task.Body, &mqTask)
			log.Printf("[Debug] starting task %v\n", mqTask)
			if err := downloadFile(mqTask, mqTask.Key); err != nil {
				log.Printf("Error while downloading %s %s\n", mqTask.Url, err)
			}
			log.Println("[Debug] downloaded file")

			go func() {
				err := c.Transcoder.Transcode(mqTask)
				if err != nil {
					log.Printf("Error while transcoding : %s\n", err)
					return
				}
				log.Println("[Debug] Transcoded finish")

				task.Ack(false)
			}()
		}
	}()
	<-forever
	return nil
}

func downloadFile(task common.MqTask, outputPath string) error {
	res, err := http.Get(task.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("status is %s", res.Status)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}
