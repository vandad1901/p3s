package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	r.GET("/users", func(ctx *gin.Context) {
		ctx.Status(200)
	})
}
