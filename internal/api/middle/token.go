package middle

import (
	"fly/pkg/httpcode"
	"fly/pkg/jwt"

	"github.com/kataras/iris/v12"
)

// AuthMiddle 验证中间件.
func AuthMiddle(ctx iris.Context) {
	token := ctx.GetHeader(jwt.Authorization)
	if token == "" {
		r, _ := httpcode.NewRequest(ctx, nil)
		r.Code(httpcode.TokenNotFound, nil)
		return
	}
	userId, err := jwt.ParseToken(token)
	if err != nil {
		r, _ := httpcode.NewRequest(ctx, nil)
		r.Code(httpcode.TokenNotValid, nil)
		return
	}
	// TODO 单点登录增加token校验
	ctx.Values().Set(jwt.CtxUserId, userId)
	ctx.Next()
}
