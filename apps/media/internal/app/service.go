package app

import (
	"github.com/vandad1901/p3s/apps/media/internal/media"
)

func initializeServices(a *App) {
	a.mediaService = media.NewService(a.logger, a.db, a.s3Client, a.mediaConsumer)
}
