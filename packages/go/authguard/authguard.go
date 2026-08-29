package authguard

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/vandad1901/p3s/packages/go/usercontext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GRPCAuthGuard(logger *slog.Logger, parser *jwt.Parser, k keyfunc.Keyfunc) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrInvalidAuth
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, ErrInvalidAuth
		}

		userID, err := HandleToken(ctx, authHeader[0], parser, k)
		if err != nil {
			logger.ErrorContext(ctx, "failed to handle token",
				"error", err,
				"method", info.FullMethod)

			return nil, ErrInvalidAuth
		}

		newCtx := usercontext.CtxWithUser(ctx, userID)

		return handler(newCtx, req)
	}
}

func EchoAuthGuard(logger *slog.Logger, parser *jwt.Parser, k keyfunc.Keyfunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return ErrInvalidAuth
			}

			userID, err := HandleToken(c.Request().Context(), tokenString, parser, k)
			if err != nil {
				logger.ErrorContext(c.Request().Context(), "failed to handle token",
					"error", err,
					"path", c.Path())

				return ErrInvalidAuth
			}

			newCtx := usercontext.CtxWithUser(c.Request().Context(), userID)

			c.SetRequest(c.Request().WithContext(newCtx))

			return next(c)
		}
	}
}

func HandleToken(ctx context.Context, tokenString string, parser *jwt.Parser, k keyfunc.Keyfunc) (int64, error) {
	var internalClaims Claims

	parsedJWT, err := parser.ParseWithClaims(tokenString, &internalClaims, k.Keyfunc)
	if err != nil {
		return 0, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if !parsedJWT.Valid {
		return 0, errInvalidJWT
	}

	id, err := strconv.ParseInt(internalClaims.Subject, 10, 64)
	if err != nil {
		return 0, errInvalidID
	}

	return id, nil
}
