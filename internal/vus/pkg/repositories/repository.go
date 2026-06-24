package repositories

import (
	"vecss/internal/storage"
	storageaws "vecss/internal/storage/aws"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type RepositoryFactory struct{}

func (f *RepositoryFactory) NewStorageRepository() storage.Storage {
	s3client := storageaws.S3Repository{
		S3Client: s3.NewFromConfig(*storageaws.AwsConfig(), func(o *s3.Options) {
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
