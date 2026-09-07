package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vandad1901/p3s/apps/upload/internal/config"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/wagslane/go-rabbitmq"
)

func initializeDependencies(cfg *config.Config) (*App, error) {
	a := new(App)
	a.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := initializeS3(a, cfg)
	if err != nil {
		return nil, err
	}

	err = initializeDatabase(a, cfg)
	if err != nil {
		return nil, err
	}

	err = initializeJWT(a, cfg)
	if err != nil {
		return nil, err
	}

	err = initializeRabbitMQ(a, cfg)
	if err != nil {
		return nil, err
	}

	return a, nil
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
	db := dbpattern.OpenDatabaseConnection(cfg.DSN)

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	a.db = db

	return nil
}

func initializeJWT(a *App, cfg *config.Config) error {
	var err error

	a.keyfunc, err = keyfunc.NewDefault([]string{cfg.AuthServiceAddress})
	if err != nil {
		return fmt.Errorf("failed to create keyfunc: %w", err)
	}

	a.parser = jwt.NewParser(
		jwt.WithValidMethods([]string{"ES256"}),
		jwt.WithIssuer("auth-service"),
		jwt.WithExpirationRequired(),
	)

	return nil
}

func initializeRabbitMQ(a *App, cfg *config.Config) error {
	rmqConn, err := rabbitmq.NewConn(cfg.RabbitMQAddress)
	if err != nil {
		return fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	a.publisher, err = rabbitmq.NewPublisher(rmqConn, rabbitmq.WithPublisherOptionsExchangeDurable)
	if err != nil {
		return fmt.Errorf("create publisher: %w", err)
	}

	return nil
}
