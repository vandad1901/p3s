package app

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/apps/media/internal/config"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/gormslog"
	"github.com/vandad1901/p3s/packages/go/rmqslog"
	"github.com/wagslane/go-rabbitmq"
)

func initializeDependencies(a *App, cfg *config.Config) error {
	err := initializeS3(a, cfg)
	if err != nil {
		return err
	}

	err = initializeDatabase(a, cfg)
	if err != nil {
		return err
	}

	err = initializeRabbitMQ(a, cfg)
	if err != nil {
		return err
	}

	return nil
}

func initializeS3(a *App, cfg *config.Config) error {
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.S3RootUsername,
				cfg.S3RootPassword,
				"",
			),
		),
	)
	if err != nil {
		return fmt.Errorf("load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		o.UsePathStyle = true
	})

	_, err = s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("connect to S3: %w", err)
	}

	a.s3Client = s3Client

	return nil
}

func initializeDatabase(a *App, cfg *config.Config) error {
	var err error

	a.db, err = dbpattern.OpenDatabaseConnection(cfg.DSN, gormslog.New(a.logger))
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}

	sqlDB, err := a.db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

func initializeRabbitMQ(a *App, cfg *config.Config) error {
	var err error

	a.rmqConn, err = rabbitmq.NewConn(cfg.RabbitMQAddress)
	if err != nil {
		return fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	a.mediaConsumer, err = rabbitmq.NewConsumer(a.rmqConn, "media",
		rabbitmq.WithConsumerOptionsConsumerAutoAck(false),
		rabbitmq.WithConsumerOptionsLogger(rmqslog.New(a.logger)),
	)
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	return nil
}
