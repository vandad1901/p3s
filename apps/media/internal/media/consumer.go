package media

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/wagslane/go-rabbitmq"
)

const consumeTimeout = 3 * time.Second

func (s *Service) RunLoop() error {
	err := s.consumer.Run(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			ctx, cancel := context.WithTimeout(context.Background(), consumeTimeout)
			defer cancel()

			err := s.handleMedia(ctx, string(d.Body))
			if err != nil {
				s.logger.Error("handling media",
					"key", string(d.Body),
					"error", err)

				if errors.Is(err, errAlreadyProcessing) || errors.Is(err, errAlreadyProcessed) {
					return rabbitmq.Ack
				}

				if errors.Is(err, errIngestFailed) {
					return rabbitmq.NackDiscard
				}

				return rabbitmq.NackRequeue
			}

			s.logger.Info("media handled successfully",
				"key", string(d.Body))

			return rabbitmq.Ack
		},
	)
	if err != nil {
		return fmt.Errorf("running media consumer loop: %w", err)
	}

	return nil
}

func (s *Service) GracefulShutdown(ctx context.Context) {
	s.consumer.CloseWithContext(ctx)
}

func (s *Service) handleMedia(ctx context.Context, key string) error {
	queuedAt, err := s.takeLease(ctx, key)
	if err != nil {
		return fmt.Errorf("creating media: %w", err)
	}

	res, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("p3s-upload-bucket"),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("getting media from S3: %w", err)
	}

	derivatives, err := IngestMedia(ctx, res.Body, *res.ContentLength)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("closing media body: %w", err)
	}

	for _, derivative := range derivatives {
		_, err = derivative.file.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("seeking derivative file: %w", err)
		}

		_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String("p3s-upload-bucket"),
			Key:    aws.String(fmt.Sprintf("%s.%s", key, derivative.Extension)),
			Body:   derivative.file,
		})
		if err != nil {
			return fmt.Errorf("putting derivative to S3: %w", err)
		}
	}

	_, err = s.changeStatus(ctx, key, queuedAt, MediaIngestStatusIngested)
	if err != nil {
		return fmt.Errorf("changing media status: %w", err)
	}

	return nil
}
