package monitor

import (
	"context"
	"sync"
	"time"
)

type TickerHandle struct {
	ctx context.Context
	wg  *sync.WaitGroup
}

// NewTickerHandle 实例化定时器处理.
func NewTickerHandle(ctx context.Context, wg *sync.WaitGroup) *TickerHandle {
	return &TickerHandle{ctx: ctx, wg: wg}
}

// Run 启动定时器执行.
func (slf *TickerHandle) Run() {
	slf.wg.Add(1)
	defer slf.wg.Done()
	ticker := time.NewTicker(time.Minute * 1)
	for range ticker.C {
		select {
		case <-slf.ctx.Done():
			return
		default:
			slf.handle()
		}
	}
}

// handle 逻辑处理.
func (slf *TickerHandle) handle() {

}
