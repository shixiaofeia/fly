package api

import (
	"fly/pkg/httpcode"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/middleware/recover"
)

// Index
func Index(app *iris.Application) {
	app.Use(recover.New())
	app.Use(httpcode.HeaderMiddleware)

	app.Options("/*", func(ctx iris.Context) {
		ctx.Next()
	})
	app.Get("/", func(ctx iris.Context) {
		r, _ := httpcode.NewRequest(ctx, nil)
		r.Ok("Welcome To Fly")
	})
	// 记载主路由
	app.Any("/debug/pprof", pprof.New())
	// 加载子路由
	app.Any("/debug/pprof/{action:path}", pprof.New())
}
