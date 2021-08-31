package config

import "fly/pkg/mysql"

type MySqlConf struct {
	Write mysql.Conf // 写配置
	Read  mysql.Conf // 读配置
}

type QiNiuConf struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
}

type ALiYunConf struct {
	AccessKeyId     string
	AccessKeySecret string
}

type WechatConf struct {
	AppId     string
	AppSecret string
}

type BaiDuConf struct {
	Ak string // 百度地图AK
}
