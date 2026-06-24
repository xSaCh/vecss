package main_test

import (
	"context"
	"testing"

	"vecss/internal/vus/pkg/repositories"

	"github.com/stretchr/testify/assert"
)

func Test_multipartUpload(t *testing.T) {
	fac := repositories.RepositoryFactory{}
	storage := fac.NewStorageRepository()

	assert.NotNil(t, storage)
	urls, err := storage.GenerateMultiPartPreSignedUrls(context.TODO(), "video.mp4", []int{1, 2, 3, 4, 5})

	assert.Nil(t, err)
	assert.Equal(t, 5, len(urls.Urls))
}
