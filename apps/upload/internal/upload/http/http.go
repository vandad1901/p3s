package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vandad1901/p3s/apps/upload/internal/upload"
)

type UploadHandler struct {
	uploadService *upload.Service
}

const maxFileSize = 20 << 20

func Register(e *echo.Group, uploadService *upload.Service) {
	handler := &UploadHandler{
		uploadService: uploadService,
	}

	e.POST("/", handler.UploadFile, middleware.BodyLimit("20M"))
}

func (h *UploadHandler) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fmt.Errorf("file required: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer func() {
		err := src.Close()
		if err != nil {
			c.Logger().Error("failed to close uploaded file")
		}
	}()

	key := c.FormValue("key")
	if key == "" {
		return errMissingKey
	}

	limitedReader := io.LimitReader(src, maxFileSize+1)

	err = h.uploadService.UploadFile(c.Request().Context(), key, limitedReader, file.Size)
	if err != nil {
		return fmt.Errorf("uploading file: %w", err)
	}

	err = c.NoContent(http.StatusOK)
	if err != nil {
		return fmt.Errorf("sending response: %w", err)
	}

	return nil
}
