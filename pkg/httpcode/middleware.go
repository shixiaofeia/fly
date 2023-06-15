package httpcode

import (
	"time"

	"github.com/kataras/iris/v12"
)

// HeaderMiddleware 设置请求头.
func HeaderMiddleware(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Authorization")
	// Access-Control-Allow-Headers 最好明确规定, 部分浏览器不兼容通配符, 参阅下文
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	ctx.Header("Access-Control-Allow-Headers", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
	if ctx.Request().Method == "OPTIONS" {
		ctx.StatusCode(200)
		return
	}
	ctx.Values().Set(CtxStartTime, time.Now())
	ctx.Next()
}
