package upload

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/packages/go/usercontext"
)

type Service struct {
	client *s3.Client
}

func NewService(client *s3.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) UploadFile(ctx context.Context, key string, file []byte) error {
	bucketName := "p3s-uploads" // hardcoded bucket name

	userID, err := usercontext.CtxUser(ctx)
	if err != nil {
		return err
	}

	uniqueKey := fmt.Sprintf("%d/%s", userID, key)

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &uniqueKey,
		Body:   bytes.NewReader(file),
	})
	if err != nil {
		return fmt.Errorf("upload to S3: %w", err)
	}

	return nil
}
