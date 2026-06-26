package main

import (
	"vecss/internal/mq"
	storageaws "vecss/internal/storage/aws"

	"vecss/internal/vus/pkg"
)

func main() {
	emitter, err := mq.NewRabbitMqEmitter("guest", "guest", "localhost")

	if err != nil {
		panic(err)
	}
	defer emitter.Connection.Close()
	emitter.Setup()

	s := storageaws.NewS3Repository()
	s.HandleBucket()

	server := pkg.NewAPIServer(":8080", s, emitter)
	server.Run()

}
