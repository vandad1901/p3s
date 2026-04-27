package main

import (
	"net/http"
	"purpl3shadow/internal/user"
	_ "purpl3shadow/utils/envutil"
	"purpl3shadow/utils/ginutil"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, ginutil.Response{
			Success: true,
			Data:    "OK",
		})
	})

	user.RegisterRoutes(r)

	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}

}
