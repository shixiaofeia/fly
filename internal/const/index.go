package constants

import "time"

// 默认配置
const (
	CallInnerTimeOut = 10 * time.Second // 请求超时
	MaxPage          = 100              // 最大请求页数
	MaxSize          = 100              // 最大请求条数
	JwtSecretKey     = "WelcomeToFly"
	Authorization    = "Authorization"
	CtxUserId        = "userId" // 用户id
)

// 通用枚举 1是2否
const (
	ModelYes = 1
	ModelNo  = 2
)
