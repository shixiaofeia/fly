package httpcode

type ErrCode struct {
	Code int
	Msg  string
}

// UpMsg 自定义错误返回.
func (c ErrCode) UpMsg(msg string) ErrCode {
	c.Msg = msg
	return c
}

var (
	Code200       = ErrCode{Code: 200, Msg: "Success"}
	ParamErr      = ErrCode{Code: 400, Msg: "参数异常"}
	TokenNotValid = ErrCode{Code: 401, Msg: "登陆已失效"}
	TokenNotFound = ErrCode{Code: 403, Msg: "Token丢失"}
	TooManyReq    = ErrCode{Code: 429, Msg: "超出请求速率"}
	ServiceErr    = ErrCode{Code: 500, Msg: "服务内部异常"}
)
