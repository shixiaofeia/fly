package monitor

import (
	"context"
	"fly/internal/config"
	"sync"
)

// Start 启动定时监控服务.
func Start(ctx context.Context, wg *sync.WaitGroup) {
	if !config.Config.IsMonitor {
		return
	}
}
