package config

import (
	"flag"
	"fly/pkg/email"
	"fly/pkg/logging"
	"fly/pkg/mongo"
	"fly/pkg/mq"
	"fly/pkg/redis"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 用来读取json 映射的结构
type jsonConfig struct {
	ServerAddress string // 对外服务地址
	ServerPort    string // 对外服务端口
	IsMonitor     bool   // 是否启动monitor
	Mysql         MySqlConf
	Mongo         mongo.Conf
	Redis         redis.Conf
	RabbitMq      mq.Conf
	Email         email.Conf
	QiNiu         QiNiuConf
	ALi           ALiYunConf
	Wechat        WechatConf
	BaiDu         BaiDuConf
}

var configPath string
var Config = jsonConfig{}

// Init 初始化函数
func Init() {
	flag.StringVar(&configPath, "config", "./config", "配置文件路径以及文件名(必填)")
	// 初始化日志
	logging.Init("./logs/fly.log", 30, true)
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	parseConfig()
	go WatchConfig()
}

// parseConfig 解析配置
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

// WatchConfig 热监听
func WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		logging.Log.Info("Conf file changed: " + e.Name)
		parseConfig()
	})
}

// ShowConfig 展示服务器运行参数
func ShowConfig() {
	logging.Log.Info("-------------------------------------------------------")
	logging.Log.Info("   服务地址:           " + Config.ServerAddress)
	logging.Log.Info("   服务端口:           " + Config.ServerPort)
	logging.Log.Infof("   健康监测:           http://%s:%s", Config.ServerAddress, Config.ServerPort)
	logging.Log.Infof("   服务监控:           http://%s:%s/debug/pprof", Config.ServerAddress, Config.ServerPort)
	logging.Log.Info("-------------------------------------------------------")
}
