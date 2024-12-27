package aws

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go"
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

func AwsReturnError(err *error) error {
	var ae smithy.APIError
	if errors.As(*err, &ae) {
		// log.Printf("[Error] %v\n", ae)
		return fmt.Errorf("%s", ae.ErrorMessage())
	}
	return *err

}
