package api

import (
	"fly/internal/httpcode"
	"github.com/gin-gonic/gin"
)

// Index router.
func Index(app *gin.Engine) {
	app.Use(httpcode.HeaderMiddleware)

	app.GET("/", func(ctx *gin.Context) {
		r, _ := httpcode.NewRequest(ctx, nil)
		r.Ok("Welcome To Fly")
	})
}
