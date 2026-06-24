package aws

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
}

func (repo *S3Repository) PutObject(ctx context.Context, filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	_, err = repo.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String(filePath),
		Body:   io.ReadSeeker(file),
	})
	if err != nil {
		return AwsReturnError(&err)
	}
	return nil
}
