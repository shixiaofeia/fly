package middle

import (
	"fly/internal/httpcode"
	"fly/pkg/redis"
	"fmt"
	"github.com/go-redis/redis_rate"
	"github.com/kataras/iris/v12"
	"golang.org/x/time/rate"
	"time"
)

// RateLimitConf 速率配置, 允许多长时间通过多少次.
type RateLimitConf struct {
	Limit int64
	Timer time.Duration
}

// exampleLimiterMap 接口请求速率配置, 建议放入redis/数据库同步本地缓存.
var exampleLimiterMap = map[string]RateLimitConf{
	"/v1/hello": {Limit: 2, Timer: time.Minute},
}

// exampleStandAloneLimiterMap 单机接口请求速率配置.
var exampleStandAloneLimiterMap = map[string]*rate.Limiter{
	"/v1/hello": rate.NewLimiter(rate.Every(time.Minute), 2),
}

// LimiterMiddle 分布式限流中间件.
func LimiterMiddle(ctx iris.Context) {
	var (
		uri    = ctx.Request().RequestURI
		client = redis.NewClusterClient()
		key    = uri
	)
	conf, ok := exampleLimiterMap[uri]
	if ok {
		limiter := redis_rate.NewLimiter(client)
		if _, _, b := limiter.Allow(key, conf.Limit, conf.Timer); !b {
			r, _ := httpcode.NewRequest(ctx, nil)
			r.Code(httpcode.TooManyReq, fmt.Errorf("req rate limit"), nil)
			return
		}
	}

	ctx.Next()
}

// StandAloneLimiterMiddle 单机限流中间件.
func StandAloneLimiterMiddle(ctx iris.Context) {
	var (
		uri = ctx.Request().RequestURI
	)
	limiter, ok := exampleStandAloneLimiterMap[uri]
	if ok {
		if b := limiter.Allow(); !b {
			r, _ := httpcode.NewRequest(ctx, nil)
			r.Code(httpcode.TooManyReq, fmt.Errorf("req rate limit"), nil)
			return
		}
	}

	ctx.Next()
}
