package app

import (
	"github.com/vandad1901/p3s/apps/media/internal/config"
	"github.com/vandad1901/p3s/apps/media/internal/media"
)

func initializeServices(cfg *config.Config, a *App) {
	a.mediaService = media.NewService(cfg.S3ExternalEndpoint, a.db, a.s3Client)
}
