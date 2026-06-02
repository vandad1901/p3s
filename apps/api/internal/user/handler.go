package user

import (
	"net/http"
	"purpl3shadow/gen/userpb"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

func getHandler(userService UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, err := strconv.Atoi(id)
		if err != nil {
			ctx.Error(err)

			return
		}

		res, err := userService.GetUserByID(ctx.Request.Context(), int64(userID))
		if err != nil {
			ctx.Error(err)

			return
		}

		json, err := protojson.Marshal(res)
		if err != nil {
			ctx.Error(err)

			return
		}

		ctx.JSON(http.StatusOK, json)
	}
}

func createHandler(userService UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := &userpb.User{}
		if err := ctx.BindJSON(user); err != nil {
			ctx.Error(err)

			return
		}

		res, err := userService.CreateUser(ctx.Request.Context(), user)
		if err != nil {
			ctx.Error(err)

			return
		}

		json, err := protojson.Marshal(res)
		if err != nil {
			ctx.Error(err)

			return
		}

		ctx.JSON(http.StatusOK, json)
	}
}
