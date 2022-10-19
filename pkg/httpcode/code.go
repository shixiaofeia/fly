package httpcode

type ErrCode struct {
	Code int
	Msg  string
}

// NewErrCode 实例化错误码
func NewErrCode(code int, msg string) ErrCode {
	return ErrCode{code, msg}
}

// UpMsg 自定义错误返回.
func (c ErrCode) UpMsg(msg string) ErrCode {
	return NewErrCode(c.Code, msg)
}

var (
	Code200       = NewErrCode(200, "Success")
	ParamErr      = NewErrCode(400, "参数异常")
	TokenNotValid = NewErrCode(401, "登陆已失效")
	TokenNotFound = NewErrCode(403, "Token丢失")
	TooManyReq    = NewErrCode(429, "超出请求速率")
	ServiceErr    = NewErrCode(500, "服务内部异常")
)
