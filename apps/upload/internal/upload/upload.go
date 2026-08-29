package upload

import (
	"bytes"
	"context"
	"fmt"
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

func (s *Service) UploadFile(ctx context.Context, key string, file []byte) error {
	userID, err := usercontext.CtxUser(ctx)
	if err != nil {
		return err
	}

	uniqueKey := fmt.Sprintf("%d/%s", userID, key)

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &uniqueKey,
		Body:        bytes.NewReader(file),
		IfNoneMatch: aws.String("*"),
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
