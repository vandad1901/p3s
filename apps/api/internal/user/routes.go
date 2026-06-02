package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userService UserService) {
	userGroup := r.Group("/users")

	userGroup.GET("/:id", getHandler(userService))
	userGroup.POST("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusCreated)
	})
}
