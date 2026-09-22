package upload

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/apps/upload/internal/outbox"
	"github.com/vandad1901/p3s/packages/go/apperror"
	"github.com/vandad1901/p3s/packages/go/usercontext"
)

type Service struct {
	s3Client        *s3.Client
	s3PresignClient *s3.PresignClient

	outboxService *outbox.Service
}

const (
	bucketName          = "p3s-upload-bucket"
	maxFileSize         = 20 << 20
	fileUploadTimeLimit = 5 * time.Minute
)

func NewService(s3Client *s3.Client, s3PresignClient *s3.PresignClient,
	outboxService *outbox.Service) *Service {
	return &Service{
		s3Client:        s3Client,
		s3PresignClient: s3PresignClient,

		outboxService: outboxService,
	}
}

func (s *Service) GenerateURL(ctx context.Context, mediaKey string) (string, map[string]string, error) {
	userID, err := usercontext.CtxUser(ctx)
	if err != nil {
		return "", nil, err
	}

	err = validateMediaKey(mediaKey)
	if err != nil {
		return "", nil, err
	}

	key := fmt.Sprintf("%d/%s", userID, mediaKey)

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    &key,
	}

	resp, err := s.s3PresignClient.PresignPostObject(ctx, params, func(o *s3.PresignPostOptions) {
		o.Expires = fileUploadTimeLimit
		o.Conditions = []any{
			[]any{"content-length-range", 0, maxFileSize},
		}
	})
	if err != nil {
		return "", nil, fmt.Errorf("generating upload URL: %w", err)
	}

	return resp.URL, resp.Values, nil
}

func (s *Service) FinalizeUpload(ctx context.Context, mediaKey string) error {
	userID, err := usercontext.CtxUser(ctx)
	if err != nil {
		return err
	}

	err = validateMediaKey(mediaKey)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d/%s", userID, mediaKey)

	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    &key,
	}

	res, err := s.s3Client.HeadObject(ctx, params)
	if err != nil {
		return fmt.Errorf("checking if file with key %s exists: %w", key, err)
	}

	if *res.ContentLength == 0 {
		return apperror.NotFound("upload.finalize.mediaNotFound")
	}

	err = s.outboxService.Enqueue(ctx, &outbox.Message{
		RoutingKey: "media",

		MessageBody: []byte(key),
	})
	if err != nil {
		return fmt.Errorf("enqueueing message: %w", err)
	}

	return nil
}
