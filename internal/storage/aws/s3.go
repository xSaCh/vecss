package aws

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"vecss/internal/domain"
)

type S3Repository struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
}

func NewS3Repository() *S3Repository {
	s3client := S3Repository{
		S3Client: s3.NewFromConfig(*AwsConfig(), func(o *s3.Options) {
			o.UsePathStyle = true
		}),
	}
	s3client.PresignClient = s3.NewPresignClient(s3client.S3Client)
	return &s3client
}

func (repo *S3Repository) PutObject(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	_, err = repo.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String(filePath),
		Body:   io.ReadSeeker(file),
	})
	if err != nil {
		return AwsReturnError(&err)
	}
	return nil
}

func (repo *S3Repository) GenerateMultiPartPreSignedUrls(ctx context.Context, key string, part []int) (*domain.MultiPartUrls, error) {
	res, err := repo.S3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String(key),
	})
	if err != nil {
		rerr := AwsReturnError(&err)
		log.Printf("[Error] %v\n", rerr)
		return nil, rerr
	}

	output := domain.MultiPartUrls{
		UploadId: *res.UploadId,
		CreateAt: time.Now(),
		ExpireAt: time.Now().Add(AWS_PRESIGN_EXPIRATION_MINTUES * time.Minute),
	}

	log.Printf("UploadId: %s\n", output.UploadId)

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

func (repo *S3Repository) CombineMultiPartUploads(ctx context.Context, input domain.CompleteMultiPartUpload) error {
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

func (repo *S3Repository) GetObjecPresignedUrl(ctx context.Context, key string) (string, error) {
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
