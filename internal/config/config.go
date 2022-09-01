package config

import (
	"fly/pkg/logging"
	"fly/pkg/safego/safe"
	"log"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Config = JsonConfig{}
)

// Init 初始化函数.
func Init(configPath string) {
	viper.SetConfigName(path.Base(configPath))
	paths := strings.Split(configPath, ".")
	if len(paths) == 0 {
		log.Fatalln("conf path err")
		return
	}
	viper.SetConfigType(paths[len(paths)-1])
	viper.AddConfigPath(path.Dir(configPath))
	parseConfig()
	safe.Go(WatchConfig)
}

// parseConfig 解析配置.
func parseConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Init config err: " + err.Error())
	}
	if err = viper.Unmarshal(&Config); err != nil {
		log.Fatalln("Unmarshal config err: " + err.Error())
	}

	logging.Init(Config.Log)
	ShowConfig()
}

// WatchConfig 热监听.
func WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		logging.Info("Conf file changed: " + e.Name)
		parseConfig()
	})
}

// ShowConfig 展示服务器运行参数.
func ShowConfig() {
	logging.Info("-------------------------------------------------------")
	logging.Info("   服务地址:           " + Config.Addr)
	logging.Info("   WEB端口:           " + Config.Port)
	logging.Info("   RPC端口:           " + Config.RpcPort)
	logging.Info("-------------------------------------------------------")
}
