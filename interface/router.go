package api

import (
	"fly/interface/example"
	"fly/pkg/httpcode"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
)

// Index
func Index(app *iris.Application) {
	app.Use(httpcode.HeaderMiddleware)

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
