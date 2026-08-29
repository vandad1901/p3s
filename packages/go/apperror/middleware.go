package apperror

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func GRPCMiddleware(logger *slog.Logger) func(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		var (
			code = codes.Internal
			msg  string
		)

		if appErr, ok := errors.AsType[Error](err); ok {
			msg = appErr.Error()

			switch appErr.Category {
			case CategoryUnauthenticated:
				code = codes.Unauthenticated
			case CategoryNotFound:
				code = codes.NotFound
			case CategoryInvalidArgument:
				code = codes.InvalidArgument
			case CategoryConflict:
				code = codes.AlreadyExists
			}
		} else {
			logger.ErrorContext(ctx, "internal error", "error", err)

			msg = uuid.NewString()
		}

		return nil, status.Error(code, msg)
	}
}

func EchoMiddleware(logger *slog.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var (
				code = http.StatusInternalServerError
				msg  string
			)

			if appErr, ok := errors.AsType[Error](err); ok {
				msg = appErr.Error()

				switch appErr.Category {
				case CategoryUnauthenticated:
					code = http.StatusUnauthorized
				case CategoryNotFound:
					code = http.StatusNotFound
				case CategoryInvalidArgument:
					code = http.StatusBadRequest
				case CategoryConflict:
					code = http.StatusConflict
				}
			} else {
				logger.ErrorContext(c.Request().Context(), "internal error", "error", err)

				msg = uuid.NewString()
			}

			return c.JSON(code, ErrorResponse{
				Message: msg,
			})
		}
	}
}
