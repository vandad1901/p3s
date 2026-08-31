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
	"gorm.io/gorm"
)

func initializeDependencies(cfg *config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	s3Client, err := initializeS3(cfg)
	if err != nil {
		return nil, err
	}

	db, err := initializeDatabase(cfg)
	if err != nil {
		return nil, err
	}

	k, err := keyfunc.NewDefault([]string{cfg.AuthServiceAddress})
	if err != nil {
		return nil, fmt.Errorf("failed to create keyfunc: %w", err)
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"ES256"}),
		jwt.WithIssuer("auth-service"),
		jwt.WithExpirationRequired(),
	)

	return &App{
		logger:   logger,
		s3Client: s3Client,
		db:       db,
		keyfunc:  k,
		parser:   parser,
	}, nil
}

func initializeS3(cfg *config.Config) (*s3.Client, error) {
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
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		o.UsePathStyle = true
	})

	_, err = s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("connect to S3: %w", err)
	}

	return s3Client, nil
}

func initializeDatabase(cfg *config.Config) (*gorm.DB, error) {
	db := dbpattern.OpenDatabaseConnection(cfg.DSN)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
