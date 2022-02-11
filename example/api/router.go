package main

import (
	"fly/example/api/example"
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
		r.JsonOk("Welcome To Fly")
	})
	// 记载主路由
	app.Any("/debug/pprof", pprof.New())
	// 加载子路由
	app.Any("/debug/pprof/{action:path}", pprof.New())
	v1 := app.Party("/v1")
	{
		example.InitApi(v1)
	}
}
