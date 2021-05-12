package middle

import (
	constants "fly/internal/const"
	"fly/pkg/httpcode"
	"github.com/kataras/iris/v12"
)

// AuthMiddle 验证中间件
func AuthMiddle(ctx iris.Context) {
	token := ctx.GetHeader(constants.Authorization)
	if token == "" {
		r, _ := httpcode.NewRequest(ctx, nil)
		r.JsonCode(httpcode.TokenNotFound, nil)
		return
	}
	userId, err := httpcode.ParseToken(token)
	if err != nil {
		r, _ := httpcode.NewRequest(ctx, nil)
		r.JsonCode(httpcode.TokenNotValid, nil)
		return
	}
	// TODO 单点登录增加token校验
	ctx.Values().Set(constants.CtxUserId, userId)
	ctx.Next()
}

// AuthOtherMiddle 验证其他中间件(有token解析, 没token就算了)
func AuthOtherMiddle(ctx iris.Context) {
	token := ctx.GetHeader(constants.Authorization)
	if token != "" {
		userId, _ := httpcode.ParseToken(token)
		ctx.Values().Set(constants.CtxUserId, userId)
	}
	ctx.Next()
}
