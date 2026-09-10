package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vandad1901/p3s/apps/upload/internal/config"
	uploadhttp "github.com/vandad1901/p3s/apps/upload/internal/upload/http"
	"github.com/vandad1901/p3s/packages/go/apperror"
	"github.com/vandad1901/p3s/packages/go/authguard"
)

func initializeServers(a *App, _ *config.Config) {
	a.echo = echo.New()

	a.echo.Pre(middleware.AddTrailingSlash())

	g := a.echo.Group("/upload/v1",
		apperror.EchoMiddleware(a.logger),
		authguard.EchoAuthGuard(a.logger, a.parser, a.keyfunc),
	)

	registerHTTPHandlers(a, g)
}

func registerHTTPHandlers(a *App, g *echo.Group) {
	uploadhttp.Register(g.Group("/upload"), a.uploadService)
}
