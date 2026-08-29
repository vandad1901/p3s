package app

import (
	"github.com/vandad1901/p3s/apps/upload/internal/config"
	"github.com/vandad1901/p3s/apps/upload/internal/upload"
)

func initializeV1Services(a *App, cfg *config.Config) {
	a.uploadService = upload.NewService(a.s3Client, cfg.S3BucketName)
}
