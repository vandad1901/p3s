package apperror

import (
	"context"
	"errors"
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

func GRPCMiddleware() func(
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

		if appErr, ok := errors.AsType[Error](err); ok {
			switch appErr.Category {
			case CategoryUnauthenticated:
				return nil, status.Error(codes.Unauthenticated, appErr.Error())
			case CategoryNotFound:
				return nil, status.Error(codes.NotFound, appErr.Error())
			case CategoryInvalidArgument:
				return nil, status.Error(codes.InvalidArgument, appErr.Error())
			case CategoryConflict:
				return nil, status.Error(codes.AlreadyExists, appErr.Error())
			}
		}

		return nil, status.Error(codes.Internal, uuid.New().String())
	}
}

func EchoMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			if appErr, ok := errors.AsType[Error](err); ok {
				switch appErr.Category {
				case CategoryUnauthenticated:
					return c.JSON(http.StatusUnauthorized, ErrorResponse{
						Message: appErr.Error(),
					})

				case CategoryNotFound:
					return c.JSON(http.StatusNotFound, ErrorResponse{
						Message: appErr.Error(),
					})

				case CategoryInvalidArgument:
					return c.JSON(http.StatusBadRequest, ErrorResponse{
						Message: appErr.Error(),
					})

				case CategoryConflict:
					return c.JSON(http.StatusConflict, ErrorResponse{
						Message: appErr.Error(),
					})
				}
			}

			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: uuid.New().String(),
			})
		}
	}
}
