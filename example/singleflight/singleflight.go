package singleflight

import "golang.org/x/sync/singleflight"

var (
	sg = &singleflight.Group{}
)

// NewSingleFlight 实例单程
// 主要解决缓存击穿的情况下保证同一key只会有一个访问DB的操作
func NewSingleFlight() *singleflight.Group {
	return sg
}
