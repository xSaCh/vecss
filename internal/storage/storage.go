package storage

import (
	"context"
	"vecss/internal/domain"
)

type Storage interface {
	PutObject(ctx context.Context, filePath string) error
	GenerateMultiPartPreSignedUrls(ctx context.Context, key string, part []int) (*domain.MultiPartUrls, error)
	CombineMultiPartUploads(ctx context.Context, input domain.CompleteMultiPartUpload) error
	GetObjecPresigntUrl(ctx context.Context, key string) (string, error)
	HandleBucket() error
}
