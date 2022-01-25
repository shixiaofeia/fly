package config

import (
	"fly/pkg/clickhouse"
	"fly/pkg/email"
	"fly/pkg/mongo"
	"fly/pkg/mq"
	"fly/pkg/mysql"
	"fly/pkg/redis"
)

// 用来读取json 映射的结构
type (
	JsonConfig struct {
		ServerAddress string // 对外服务地址
		ServerPort    string // 对外服务端口
		RpcPort       string // rpc端口
		IsMonitor     bool   // 是否启动monitor
		Mysql         MySqlConf
		Mongo         mongo.Conf
		Redis         redis.Conf
		RabbitMq      mq.Conf
		Es            []string
		Clickhouse    clickhouse.Config
		Email         email.Conf
		QiNiu         QiNiuConf
		ALi           ALiYunConf
		Wechat        WechatConf
		BaiDu         BaiDuConf
	}
	MySqlConf struct {
		Write mysql.Conf // 写配置
		Read  mysql.Conf // 读配置
	}
	QiNiuConf struct {
		AccessKey string
		SecretKey string
		Bucket    string
		Domain    string
	}
	ALiYunConf struct {
		AccessKeyId     string
		AccessKeySecret string
	}
	WechatConf struct {
		AppId     string
		AppSecret string
	}
	BaiDuConf struct {
		Ak string // 百度地图AK
	}
)
