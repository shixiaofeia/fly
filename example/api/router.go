package main

import (
	"fly/example/api/example"
	"fly/internal/api/middle"
	"fly/internal/httpcode"
	"github.com/gin-gonic/gin"
)

// Index 路由
func Index(app *gin.Engine) {
	app.Use(httpcode.HeaderMiddleware)

	app.GET("/", func(ctx *gin.Context) {
		r, _ := httpcode.NewRequest(ctx, nil)
		r.Ok("Welcome To Fly")
	})

	v1 := app.Group("/v1")
	v1.Use(middle.StandAloneLimiterMiddle)
	{
		example.InitApi(v1)
	}
}
