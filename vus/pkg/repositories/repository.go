package repositories

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/xSaCh/vecss/vus/pkg/repositories/aws"
)

type RepositoryFactory struct{}

func (f *RepositoryFactory) NewStorageRepository() Storage {
	s3client := aws.S3Repository{
		S3Client: s3.NewFromConfig(*aws.AwsConfig(), func(o *s3.Options) {
			o.UsePathStyle = true
		}),
	}
	s3client.PresignClient = s3.NewPresignClient(s3client.S3Client)

	err := s3client.HandleBucket()
	if err != nil {
		return nil
	}

	return &s3client
}
