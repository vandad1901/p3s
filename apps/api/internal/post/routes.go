package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func RegisterRoutes(r *gin.Engine, postService PostService) {
	r.GET("/posts/:id", GetHandler(postService))

	r.POST("/posts", func(ctx *gin.Context) {
		ctx.Status(http.StatusCreated)
	})
}

func GetHandler(postService PostService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		postID, err := strconv.Atoi(id)
		if err != nil {
			ctx.Status(http.StatusBadRequest)

			return
		}

		res, err := postService.GetPostByID(ctx, int64(postID))
		if err != nil {
			ctx.Status(http.StatusInternalServerError)

			return
		}

		json, err := proto.Marshal(res)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)

			return
		}

		ctx.JSON(http.StatusOK, json)
	}
}
