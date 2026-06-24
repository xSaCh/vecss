package repositories

import (
	"context"
	"vecss/internal/domain"
)

type Storage interface {
	// Upload(ctx context.Context, bucketName string, objectKey string, data []byte) error
	// Download(ctx context.Context, bucketName string, objectKey string) ([]byte, error)
	// Delete(ctx context.Context, bucketName string, objectKey string) error
	T()
	GenerateMultiPartPreSignedUrls(ctx context.Context, key string, part []int) (*domain.MultiPartUrls, error)
	CombineMultiPartUploads(ctx context.Context, input domain.CompleteMultiPartUpload) error
	GetObjecPresigntUrl(ctx context.Context, key string) (string, error)
}
