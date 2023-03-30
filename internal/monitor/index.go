package monitor

import (
	"context"
	"fly/internal/config"
	"fly/pkg/safego/safe"
	"sync"
)

// Start 启动定时监控服务.
func Start(ctx context.Context, wg *sync.WaitGroup) {
	if !config.Config.IsMonitor {
		return
	}

	safe.Go(NewCronHandle(ctx, wg).Run)
	safe.Go(NewTickerHandle(ctx, wg).Run)
}
