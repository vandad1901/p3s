package ginutil

import "github.com/gin-gonic/gin"

func validateProtoJSON[T any](ctx *gin.Context, obj *T) bool {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		ctx.JSON(400, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    "invalid_request",
				Message: "Invalid JSON body",
			},
		})

		return false
	}
	return true
}
