package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Repository struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
}

func (repo *S3Repository) T() {

	name := "bkt"
	o, err := repo.S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint("us-east-1"),
		},
	},
	)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			log.Printf("You already own bucket %s.\n", name)
			err = owned
		} else if errors.As(err, &exists) {
			log.Printf("Bucket %s already exists.\n", name)
			err = exists
		}
	} else {
		err = s3.NewBucketExistsWaiter(repo.S3Client).Wait(
			context.TODO(), &s3.HeadBucketInput{Bucket: aws.String(name)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for bucket %s to exist.\n", name)
		}
	}

	fmt.Printf("o: %v\n", o)
}

func (repo *S3Repository) GenerateMultiPartPreSignedUrls(ctx context.Context, part []int) ([]*v4.PresignedHTTPRequest, error) {

	uploadInput := &s3.CreateMultipartUploadInput{
		Bucket: aws.String("bkt"),
		Key:    aws.String("video.mp4"),
	}

	res, err := repo.S3Client.CreateMultipartUpload(ctx, uploadInput)
	if err != nil {
		return nil, err
	}

	fmt.Printf("res.UploadId: %s\n", *res.UploadId)

	var urls []*v4.PresignedHTTPRequest
	for _, p := range part {
		req, err := repo.PresignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String("bkt"),
			Key:        aws.String("video.mp4"),
			UploadId:   res.UploadId,
			PartNumber: aws.Int32(int32(p)),
		}, func(options *s3.PresignOptions) {
			options.Expires = time.Duration(15 * int64(time.Minute))
		})
		if err != nil {
			return nil, err
		}
		urls = append(urls, req)
	}
	return urls, nil
}

func (repo *S3Repository) CombineMultiPartUploads(ctx context.Context, uploadId string, etags []string, partnumbers []int) error {
	var parts []types.CompletedPart
	for i, etag := range etags {
		parts = append(parts, types.CompletedPart{
			ETag:       aws.String(etag),
			PartNumber: aws.Int32(int32(partnumbers[i])),
		})
	}
	_, err := repo.S3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String("bkt"),
		Key:             aws.String("video.mp4"),
		UploadId:        aws.String(uploadId),
		MultipartUpload: &types.CompletedMultipartUpload{Parts: parts},
	})
	return err
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (repo *S3Repository) GetObject(
	ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := repo.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified number of seconds.
func (repo *S3Repository) PutObject(
	ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := repo.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

// DeleteObject makes a presigned request that can be used to delete an object from a bucket.
func (repo *S3Repository) DeleteObject(ctx context.Context, bucketName string, objectKey string) (*v4.PresignedHTTPRequest, error) {
	request, err := repo.PresignClient.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to delete object %v. Here's why: %v\n", objectKey, err)
	}
	return request, err
}

func (repo *S3Repository) PresignPostObject(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*s3.PresignedPostRequest, error) {
	request, err := repo.PresignClient.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(options *s3.PresignPostOptions) {
		options.Expires = time.Duration(lifetimeSecs) * time.Second
	})
	if err != nil {
		log.Printf("Couldn't get a presigned post request to put %v:%v. Here's why: %v\n", bucketName, objectKey, err)
	}
	return request, nil
}
