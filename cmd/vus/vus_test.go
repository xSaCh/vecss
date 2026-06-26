package main_test

import (
	"context"
	"testing"

	storageaws "vecss/internal/storage/aws"

	"github.com/stretchr/testify/assert"
)

func Test_multipartUpload(t *testing.T) {
	s := storageaws.NewS3Repository()
	s.HandleBucket()

	assert.NotNil(t, s)
	urls, err := s.GenerateMultiPartPreSignedUrls(context.TODO(), "video.mp4", []int{1, 2, 3, 4, 5})

	assert.Nil(t, err)
	assert.Equal(t, 5, len(urls.Urls))
}
