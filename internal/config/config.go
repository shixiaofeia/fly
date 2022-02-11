package config

import (
	"fly/pkg/logging"
	"fly/pkg/safego/safe"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Config = JsonConfig{}
)

// Init 初始化函数.
func Init(configPath string) {

	// 初始化日志
	logging.Init("./logs/fly.log", 30, true)
	viper.SetConfigName(path.Base(configPath))
	viper.SetConfigType("json")
	viper.AddConfigPath(path.Dir(configPath))
	parseConfig()
	safe.Go(WatchConfig)
}

// parseConfig 解析配置.
func parseConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		logging.Log.Fatal("Init config err: " + err.Error())
	}
	if err = viper.Unmarshal(&Config); err != nil {
		logging.Log.Fatal("Unmarshal config err: " + err.Error())
	}
	ShowConfig()
}

// WatchConfig 热监听.
func WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		logging.Log.Info("Conf file changed: " + e.Name)
		parseConfig()
	})
}

// ShowConfig 展示服务器运行参数.
func ShowConfig() {
	logging.Log.Info("-------------------------------------------------------")
	logging.Log.Info("   服务地址:           " + Config.ServerAddress)
	logging.Log.Info("   WEB端口:           " + Config.ServerPort)
	logging.Log.Info("   RPC端口:           " + Config.RpcPort)
	logging.Log.Info("-------------------------------------------------------")
}
