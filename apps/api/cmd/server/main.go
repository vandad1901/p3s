package main

import (
	"purpl3shadow/internal/post"
	"purpl3shadow/internal/user"
	"purpl3shadow/utils/dbutil"
	_ "purpl3shadow/utils/envutil"
	"purpl3shadow/utils/ginutil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync()

	r := gin.Default()
	r.Use(ginutil.LogAttacher(logger), ginutil.ErrorHandler)

	userService := user.NewUserService(dbutil.DB)
	postService := post.NewPostService(dbutil.DB)

	user.RegisterRoutes(r, userService)
	post.RegisterRoutes(r, postService)

	err := r.Run("localhost:8080")
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
