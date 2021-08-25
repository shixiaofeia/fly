package httpcode

import (
	"github.com/kataras/iris/v12"
	"time"
)

// HeaderMiddleware 设置请求头
func HeaderMiddleware(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "*")
	ctx.Header("content-type", "application/json;charset=utf-8")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.StatusCode(204)
		return
	}
	start := time.Now().UnixNano()
	ctx.Values().Set(CtxStartTime, start)
	ctx.Next()
}
