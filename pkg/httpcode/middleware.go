package httpcode

import (
	"fly/pkg/logging"
	"github.com/kataras/iris/v12"
	"log"
	"runtime"
	"time"
)

// HeaderMiddleware 设置请求头
func HeaderMiddleware(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 打印堆栈信息
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			// 此处正常返回结果, 不打印panic日志
			logging.Log.Error("service recover err-------------------")
			log.Println(string(buf))
			r, _ := NewRequest(ctx, nil)
			r.JsonCode(ServiceErr, err)
			return
		}
	}()
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Origin", ctx.Request().Header.Get("origin"))
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,x-access-token,captcha-session")
	ctx.Header("content-type", "application/json")
	start := time.Now().UnixNano()
	ctx.Values().Set(CtxStartTime, start)
	ctx.Next()
}
