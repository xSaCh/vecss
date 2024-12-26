package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const AWS_REGION = "us-east-1"
const AWS_BUCKET = "bkt"
const AWS_PRESIGN_EXPIRATION_MINTUES = 15

func AwsConfig() *aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AWS_REGION))
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
