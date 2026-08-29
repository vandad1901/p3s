package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vandad1901/p3s/apps/auth/internal/jwks"
)

type JWKSHandler struct {
	jwksService *jwks.Service
}

func Register(e *echo.Echo, uploadService *jwks.Service) {
	handler := &JWKSHandler{
		jwksService: uploadService,
	}

	e.GET("/.well-known/jwks.json", handler.Handle)
}

func (h *JWKSHandler) Handle(c echo.Context) error {
	response, err := h.jwksService.Handle(c.Request().Context())
	if err != nil {
		return fmt.Errorf("reading keyset: %w", err)
	}

	err = c.JSON(http.StatusOK, response)
	if err != nil {
		return fmt.Errorf("sending response: %w", err)
	}

	return nil
}
