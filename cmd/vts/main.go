package main

import (
	"context"
	"log"
	"vecss/internal/mq"

	"vecss/internal/vts"

	storageaws "vecss/internal/storage/aws"
)

const NUM_WORKERS = 2

func main() {
	rbmq, err := mq.NewRabbitMqEmitter("guest", "guest", "localhost")

	if err != nil {
		panic(err)
	}
	defer rbmq.Connection.Close()
	if err := rbmq.Setup(); err != nil {
		log.Fatalf("Failed to setup rabbitmq: %v", err)
		panic(err)
	}

	if err := rbmq.Channel.Qos(
		NUM_WORKERS, // prefetch count
		0,           // prefetch size
		false,       // global
	); err != nil {
		log.Fatalf("Failed to setup Qos: %v", err)
		panic(err)
	}

	s3Repo := storageaws.NewS3Repository()

	t := vts.FFMpegTranscoder{}
	con := vts.NewConsumer(rbmq, &t, s3Repo)

	con.Listen(context.TODO(), NUM_WORKERS)

}
