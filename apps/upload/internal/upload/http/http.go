package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vandad1901/p3s/apps/upload/internal/upload"
)

type UploadHandler struct {
	uploadService *upload.Service
}

func Register(e *echo.Group, uploadService *upload.Service) {
	handler := &UploadHandler{
		uploadService: uploadService,
	}

	e.POST("/init/", handler.UploadFile)
	e.POST("/finalize/", handler.FinalizeUpload)
}

func (h *UploadHandler) UploadFile(c echo.Context) error {
	key := c.FormValue("key")
	if key == "" {
		return errMissingKey
	}

	res, err := h.uploadService.GenerateURL(c.Request().Context(), key)
	if err != nil {
		return fmt.Errorf("uploading file: %w", err)
	}

	err = c.JSON(http.StatusOK, res)
	if err != nil {
		return fmt.Errorf("sending response: %w", err)
	}

	return nil
}

func (h *UploadHandler) FinalizeUpload(c echo.Context) error {
	key := c.FormValue("key")
	if key == "" {
		return errMissingKey
	}

	err := h.uploadService.FinalizeUpload(c.Request().Context(), key)
	if err != nil {
		return fmt.Errorf("finalizing upload: %w", err)
	}

	err = c.NoContent(http.StatusOK)
	if err != nil {
		return fmt.Errorf("sending response: %w", err)
	}

	return nil
}
