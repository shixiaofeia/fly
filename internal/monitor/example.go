package monitor

import (
	"context"
	"time"

	"github.com/robfig/cron"
)

// StartCron 启动定时任务.
func StartCron() {
	c := cron.New()
	_ = c.AddFunc("0 0 0 */1 * ?", handle)
	c.Start()
}

// StartTicker
func StartTicker(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 1)
	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			handle()
		}
	}
}

// handle 逻辑处理.
func handle() {

}
