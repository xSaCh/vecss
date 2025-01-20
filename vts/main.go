package main

import (
	"common/mq"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	vts "github.com/xSaCh/vecss/vts/package"
	"github.com/xSaCh/vecss/vts/package/aws"
)

func main() {
	emitter, err := mq.NewRabbitMqEmitter("guest", "guest", "localhost")

	if err != nil {
		panic(err)
	}
	defer emitter.Connection.Close()
	emitter.Setup()

	s3client := aws.S3Repository{
		S3Client: s3.NewFromConfig(*aws.AwsConfig(), func(o *s3.Options) {
			o.UsePathStyle = true
		}),
	}

	t := vts.FFMpegTranscoder{}
	con := vts.NewConsumer(emitter, &t, &s3client)

	con.Listen(context.TODO())

}
