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

func getErrorInfoGRPC(err error) (codes.Code, string) {
	if appErr, ok := errors.AsType[Error](err); ok {
		msg := appErr.Error()

		switch appErr.Category {
		case CategoryUnauthenticated:
			return codes.Unauthenticated, msg
		case CategoryNotFound:
			return codes.NotFound, msg
		case CategoryInvalidArgument:
			return codes.InvalidArgument, msg
		case CategoryConflict:
			return codes.AlreadyExists, msg
		}
	}

	return codes.Internal, uuid.NewString()
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

		code, msg := getErrorInfoGRPC(err)

		logger.ErrorContext(ctx, "internal error", "error", err, "message", msg)

		return nil, status.Error(code, msg)
	}
}

func getErrorInfoHTTP(err error) (int, string) {
	if appErr, ok := errors.AsType[Error](err); ok {
		msg := appErr.Error()

		switch appErr.Category {
		case CategoryUnauthenticated:
			return http.StatusUnauthorized, msg
		case CategoryNotFound:
			return http.StatusNotFound, msg
		case CategoryInvalidArgument:
			return http.StatusBadRequest, msg
		case CategoryConflict:
			return http.StatusConflict, msg
		}
	}

	if httpErr, ok := errors.AsType[*echo.HTTPError](err); ok {
		if httpErr.Code == http.StatusNotFound {
			return httpErr.Code, "Not Found"
		}
	}

	return http.StatusInternalServerError, uuid.NewString()
}

func EchoMiddleware(logger *slog.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			code, msg := getErrorInfoHTTP(err)

			logger.ErrorContext(c.Request().Context(), "internal error", "error", err, "message", msg)

			return c.JSON(code, ErrorResponse{
				Message: msg,
			})
		}
	}
}
