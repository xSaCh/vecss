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

	common "common"
)

type S3Repository struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
}

func (repo *S3Repository) T() {

	ott, err := repo.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("bkt"),
		Key:    aws.String("b"),
	})
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}
	fmt.Printf("url: %v\n", ott.URL)
	return
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

func (repo *S3Repository) GenerateMultiPartPreSignedUrls(ctx context.Context, key string, part []int) (*common.MultiPartUrls, error) {

	res, err := repo.S3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String(key),
	})
	if err != nil {
		rerr := AwsReturnError(&err)
		log.Printf("[Error] %v\n", rerr)
		return nil, rerr
	}

	output := common.MultiPartUrls{
		UploadId: *res.UploadId,
		CreateAt: time.Now(),
		ExpireAt: time.Now().Add(AWS_PRESIGN_EXPIRATION_MINTUES * time.Minute),
	}

	log.Printf("UploadId: %s\n", output.UploadId)

	// Generate presigned URLs for each part
	var urls []*v4.PresignedHTTPRequest
	for _, p := range part {
		req, err := repo.PresignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(AWS_BUCKET),
			Key:        aws.String(key),
			UploadId:   res.UploadId,
			PartNumber: aws.Int32(int32(p)),
		}, func(options *s3.PresignOptions) {
			options.Expires = AWS_PRESIGN_EXPIRATION_MINTUES * time.Minute
		})
		if err != nil {
			rerr := AwsReturnError(&err)
			log.Printf("[Error] %v\n", rerr)
			return nil, rerr
		}
		urls = append(urls, req)
	}

	output.Urls = make([]string, len(urls))
	for i, u := range urls {
		output.Urls[i] = u.URL
	}
	return &output, nil
}

func (repo *S3Repository) CombineMultiPartUploads(ctx context.Context, input common.CompleteMultiPartUpload) error {
	var parts []types.CompletedPart
	for i, etag := range input.ETags {
		parts = append(parts, types.CompletedPart{
			ETag:       aws.String(etag),
			PartNumber: aws.Int32(int32(input.PartNumbers[i])),
		})
	}
	_, err := repo.S3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		UploadId:        aws.String(input.UploadId),
		Bucket:          aws.String(AWS_BUCKET),
		Key:             aws.String(input.Key),
		MultipartUpload: &types.CompletedMultipartUpload{Parts: parts},
	})

	if err != nil {
		rerr := AwsReturnError(&err)
		log.Printf("[Error] %v\n", rerr)
		return rerr
	}

	return nil
}

func (repo *S3Repository) GetObjecPresigntUrl(ctx context.Context, key string) (string, error) {

	res, err := repo.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String(key),
	})

	if err != nil {
		rerr := AwsReturnError(&err)
		log.Printf("[Error] %v\n", rerr)
		return "", rerr
	}

	return res.URL, nil
}

func (repo *S3Repository) HandleBucket() error {

	_, err := repo.S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(AWS_BUCKET),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(AWS_REGION),
		},
	},
	)
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			log.Printf("You already own bucket %s.\n", AWS_BUCKET)
			err = owned
		} else if errors.As(err, &exists) {
			log.Printf("Bucket %s already exists.\n", AWS_BUCKET)
			err = exists
		}
	} else {
		err = s3.NewBucketExistsWaiter(repo.S3Client).Wait(
			context.TODO(), &s3.HeadBucketInput{Bucket: aws.String(AWS_BUCKET)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for bucket %s to exist.\n", AWS_BUCKET)
		}
	}
	return nil
}

// Makes a presigned request that can be used to put an object in a bucket.
func (repo *S3Repository) genereatePresignObjectUrl(
	ctx context.Context, bucketName string, objectKey string, validFor time.Duration) (string, error) {
	request, err := repo.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = validFor
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}

	return request.URL, err
}
