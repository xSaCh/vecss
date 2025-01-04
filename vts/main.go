package main

import (
	"common/mq"
	"context"

	vts "github.com/xSaCh/vecss/vts/package"
)

func main() {
	emitter, err := mq.NewRabbitMqEmitter("guest", "guest", "localhost")

	if err != nil {
		panic(err)
	}
	defer emitter.Connection.Close()
	emitter.Setup()

	t := vts.FFMpegTranscoder{}
	con := vts.NewConsumer(emitter, &t)

	con.Listen(context.TODO())

}
