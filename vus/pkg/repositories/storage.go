package repositories

import (
	"common"
	"context"
)

type Storage interface {
	// Upload(ctx context.Context, bucketName string, objectKey string, data []byte) error
	// Download(ctx context.Context, bucketName string, objectKey string) ([]byte, error)
	// Delete(ctx context.Context, bucketName string, objectKey string) error
	T()
	GenerateMultiPartPreSignedUrls(ctx context.Context, key string, part []int) (*common.MultiPartUrls, error)
	CombineMultiPartUploads(ctx context.Context, input common.CompleteMultiPartUpload) error
	GetObjecPresigntUrl(ctx context.Context, key string) (string, error)
}
