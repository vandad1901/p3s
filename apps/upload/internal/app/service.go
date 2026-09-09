package app

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/apps/upload/internal/outbox"
	"github.com/vandad1901/p3s/apps/upload/internal/upload"
)

func initializeServices(a *App) {
	a.outboxService = outbox.NewService(a.logger, a.db, a.publisher)
	a.uploadService = upload.NewService(a.s3Client, s3.NewPresignClient(a.s3Client), a.outboxService)
}
