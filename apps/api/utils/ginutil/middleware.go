package ginutil

import (
	"context"
	"net/http"
	"purpl3shadow/utils/obviousutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func LogAttacher(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestLogger := logger.With(zap.String("URL", ctx.Request.URL.Path))
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), loggerKey, requestLogger))

		ctx.Next()
	}
}

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()

	logger := ctx.Request.Context().Value(loggerKey).(*zap.Logger)

	if len(ctx.Errors) > 0 {
		if echoedError := ctx.Errors[0].Err; obviousutil.IsObvious(echoedError) {
			ctx.JSON(http.StatusBadRequest,
				&Response{
					Error: &ErrorInfo{
						Key: ctx.Errors[0].Error()}})
		} else {
			errorUUID, err := uuid.NewRandom()
			if err != nil {
				errorUUID = uuid.Nil
			}

			uuidStr := errorUUID.String()

			ctx.JSON(http.StatusInternalServerError,
				&Response{
					Error: &ErrorInfo{
						Key:      uuidStr,
						Internal: true}})

			logger.Warn("Internal server error",
				zap.String("UUID", uuidStr),
				zap.Error(echoedError))
		}
	}
}
