package monitor

import (
	"context"
	"fly/config"
)

// InitMonitor 初始化定时监控服务
func InitMonitor(ctx context.Context) {
	if !config.Config.IsMonitor {
		return
	}
}
