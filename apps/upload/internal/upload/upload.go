package upload

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/packages/go/usercontext"
)

type Service struct {
	client *s3.Client

	bucketName string
}

func NewService(client *s3.Client, bucketName string) *Service {
	return &Service{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *Service) UploadFile(ctx context.Context, key string,
	fileReader io.Reader, fileSize int64) error {
	userID, err := usercontext.CtxUser(ctx)
	if err != nil {
		return err
	}

	uniqueKey := fmt.Sprintf("%d/%s", userID, key)

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &s.bucketName,
		Key:           &uniqueKey,
		Body:          fileReader,
		ContentLength: &fileSize,
		IfNoneMatch:   aws.String("*"),
		Metadata: map[string]string{
			"uploaded-by": strconv.FormatInt(userID, 10),
		},
	},
	)
	if err != nil {
		return fmt.Errorf("upload to S3: %w", err)
	}

	return nil
}

func (s *Service) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	userID, err := usercontext.CtxUser(ctx)
	if err != nil {
		return nil, err
	}

	uniqueKey := fmt.Sprintf("%d/%s", userID, key)

	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &uniqueKey,
	})
	if err != nil {
		return nil, fmt.Errorf("get from S3: %w", err)
	}

	return resp.Body, nil
}

func (s *Service) DeleteFile(ctx context.Context, key string) error {
	userID, err := usercontext.CtxUser(ctx)
	if err != nil {
		return err
	}

	uniqueKey := fmt.Sprintf("%d/%s", userID, key)

	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucketName,
		Key:    &uniqueKey,
	})
	if err != nil {
		return fmt.Errorf("delete from S3: %w", err)
	}

	return nil
}
