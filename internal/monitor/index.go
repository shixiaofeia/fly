package monitor

import "fly/config"

// InitMonitor 初始化定时监控服务
func InitMonitor() {
	if !config.Config.IsMonitor {
		return
	}
}
