package app

import (
	"github.com/labstack/echo/v4"
	"github.com/vandad1901/p3s/apps/upload/internal/config"
	uploadhttp "github.com/vandad1901/p3s/apps/upload/internal/upload/http"
	"github.com/vandad1901/p3s/packages/go/apperror"
	"github.com/vandad1901/p3s/packages/go/authguard"
)

func initializeServers(a *App, cfg *config.Config) {
	a.httpServer = echo.New()

	g := a.httpServer.Group("/v1",
		apperror.EchoMiddleware(a.logger),
		authguard.EchoAuthGuard(a.logger, a.parser, a.keyfunc))

	registerHTTPHandlers(a, g)
}

func registerHTTPHandlers(a *App, g *echo.Group) {
	uploadhttp.Register(g.Group("upload"), a.uploadService)
}
