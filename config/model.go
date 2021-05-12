package config

import "fly/pkg/mysql"

type MySqlConfig struct {
	Write mysql.Conf // 写配置
	Read  mysql.Conf // 读配置
}

type QiNiuConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
}

type ALiYunConfig struct {
	AccessKeyId     string
	AccessKeySecret string
}

type WechatConfig struct {
	AppId     string
	AppSecret string
}

type BaiDuConfig struct {
	Ak string // 百度地图AK
}
