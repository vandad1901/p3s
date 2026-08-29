package app

import "github.com/vandad1901/p3s/apps/upload/internal/upload"

func initializeV1Services(a *App) {
	a.uploadService = upload.NewService(a.s3Client)
}
