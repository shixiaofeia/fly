package main

import (
	"fly/example/api/example"
	"fly/internal/api/middle"
	"fly/pkg/httpcode"
	"github.com/kataras/iris/v12/middleware/logger"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/middleware/recover"
)

// Index 路由
func Index(app *iris.Application) {
	app.Use(recover.New())
	app.Use(httpcode.HeaderMiddleware)
	app.Use(logger.New(logger.DefaultConfig()))

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
	v1 := app.Party("/v1")
	v1.Use(middle.StandAloneLimiterMiddle)
	{
		example.InitApi(v1)
	}
}
