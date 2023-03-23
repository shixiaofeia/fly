package monitor

import (
	"context"
	"github.com/robfig/cron"
	"sync"
)

type CronHandle struct {
	ctx context.Context
	wg  *sync.WaitGroup
}

// NewCronHandle 实例定时任务处理.
func NewCronHandle(ctx context.Context, wg *sync.WaitGroup) *CronHandle {
	return &CronHandle{ctx: ctx, wg: wg}
}

// Run 启动cron.
func (slf *CronHandle) Run() {
	c := cron.New()
	_ = c.AddFunc("0 0 0 */1 * ?", slf.handle)
	c.Start()
}

// handle 逻辑处理.
func (slf *CronHandle) handle() {
	select {
	case <-slf.ctx.Done():
		return
	default:
	}
	slf.wg.Add(1)
	defer slf.wg.Done()

}
