package app

import (
	"github.com/vandad1901/p3s/apps/media/internal/config"
	"github.com/vandad1901/p3s/apps/media/internal/media"
)

func initializeServices(a *App, _ *config.Config) {
	a.mediaService = media.NewService(a.db, a.mediaConsumer)
}
