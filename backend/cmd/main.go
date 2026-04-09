package main

import (
	"purpl3shadow/backend/internal/user"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func main() {
	r := gin.Default()

	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(200, Response{
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
