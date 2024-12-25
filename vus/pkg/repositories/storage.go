package repositories

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	// Upload(ctx context.Context, bucketName string, objectKey string, data []byte) error
	// Download(ctx context.Context, bucketName string, objectKey string) ([]byte, error)
	// Delete(ctx context.Context, bucketName string, objectKey string) error
	T()
	PutObject(
		ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error)
	PresignPostObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*s3.PresignedPostRequest, error)
}
