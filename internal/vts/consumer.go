package vts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"vecss/internal/mq"
	"vecss/internal/storage"

	"vecss/internal/domain"
)

type Consumer struct {
	Emitter    mq.Emitter
	Transcoder Transcoder
	Storage    storage.Storage
}

func NewConsumer(rbmq *mq.RabbitMq, transcoder Transcoder, storage storage.Storage) *Consumer {
	return &Consumer{
		Emitter:    rbmq,
		Transcoder: transcoder,
		Storage:    storage,
	}
}

func (c *Consumer) Listen(ctx context.Context) error {
	tasks, err := c.Emitter.Consume(ctx)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for task := range tasks {
			var mqTask domain.MqTask
			json.Unmarshal(task.Body(), &mqTask)
			log.Printf("[Debug] starting task %v\n", mqTask)

			downloadedFilePath, err := downloadFile(mqTask, mqTask.Key)
			if err != nil {
				log.Printf("Error while downloading %s %s\n", mqTask.Url, err)
			}
			log.Printf("[Debug] downloaded file at %s\n", downloadedFilePath)

			go func() {
				paths, err := c.Transcoder.Transcode(mqTask, downloadedFilePath)
				if err != nil {
					log.Printf("Error while transcoding : %s\n", err)
					task.Nack(true)
					return
				}
				log.Println("[Debug] Transcoded finish")

				for _, p := range paths {

					err = c.Storage.PutObject(ctx, p)
					if err != nil {
						log.Printf("Error while uploading %s : %s\n", p, err)
						task.Nack(true)
						return
					}
				}
				log.Println("[Debug] Uploading finish")

				task.Ack()
			}()
		}
	}()
	<-forever
	return nil
}

func downloadFile(task domain.MqTask, fileKey string) (string, error) {
	res, err := http.Get(task.Url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status is %s", res.Status)
	}

	tempFile, err := os.CreateTemp("", fmt.Sprintf("%s-*", fileKey))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, res.Body)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
