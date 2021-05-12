package redis

import (
	"fly/pkg/logging"
	"fmt"
	"github.com/go-redis/redis"
)

var clusterClient *redis.ClusterClient
var standAloneClient *redis.Client

type Conf struct {
	Address     string   // 单机地址
	Password    string   // 密码
	AddressList []string // 集群地址
}

// Init 初始化redis服务
func Init(c Conf) (err error) {
	if err = InitCluster(c.AddressList, c.Password); err != nil {
		return
	}
	if err = InitStandAlone(c.Address, c.Password); err != nil {
		return
	}
	return
}

// NewClusterClient
func NewClusterClient() *redis.ClusterClient {
	return clusterClient
}

// NewStandAloneClient
func NewStandAloneClient() *redis.Client {
	return standAloneClient
}

// InitCluster  // 初始化redis集群模式
func InitCluster(address []string, passWord string) (err error) {
	if len(address) == 0 {
		logging.Log.Warn("InitCluster not config")
		return
	}
	clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    address,
		PoolSize: 1000,
		Password: passWord,
	})
	if err = clusterClient.ReloadState(); err != nil {
		return fmt.Errorf("InitCluster err: %v", err)
	}
	return
}

// InitStandAlone 初始化redis单机模式
func InitStandAlone(address string, passWord string) (err error) {
	if address == "" {
		logging.Log.Warn("InitStandAlone not config")
		return
	}
	standAloneClient = redis.NewClient(
		&redis.Options{
			Addr:         address,
			Password:     passWord,
			MinIdleConns: 100,
		})
	return
}
