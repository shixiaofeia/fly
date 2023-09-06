package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

var clusterClient *redis.ClusterClient
var standAloneClient *redis.Client

type Conf struct {
	Addr     string   // 单机地址
	Pwd      string   // 密码
	AddrList []string // 集群地址
}

// Init 初始化redis服务.
func Init(c Conf) (err error) {
	if err = InitCluster(c.AddrList, c.Pwd); err != nil {
		return
	}
	if err = InitStandAlone(c.Addr, c.Pwd); err != nil {
		return
	}
	return
}

// NewClusterClient 获取集群连接
func NewClusterClient() *redis.ClusterClient {
	return clusterClient
}

// NewStandAloneClient 获取单机连接
func NewStandAloneClient() *redis.Client {
	return standAloneClient
}

// InitCluster  // 初始化redis集群模式.
func InitCluster(address []string, passWord string) (err error) {
	if len(address) == 0 {
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

// InitStandAlone 初始化redis单机模式.
func InitStandAlone(address string, passWord string) (err error) {
	if address == "" {
		return
	}
	standAloneClient = redis.NewClient(
		&redis.Options{
			Addr:         address,
			Password:     passWord,
			MinIdleConns: 100,
		})
	if err = standAloneClient.Ping().Err(); err != nil {
		return err
	}

	return
}
