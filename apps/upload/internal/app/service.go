package app

import (
	"github.com/vandad1901/p3s/apps/upload/internal/config"
	"github.com/vandad1901/p3s/apps/upload/internal/mediaupload"
	"github.com/vandad1901/p3s/apps/upload/internal/upload"
)

func initializeServices(a *App, cfg *config.Config) {
	a.mediaUploadService = mediaupload.NewService(a.db)
	a.uploadService = upload.NewService(a.s3Client, cfg.S3BucketName, a.mediaUploadService)
}
