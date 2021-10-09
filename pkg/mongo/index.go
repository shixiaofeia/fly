package mongo

import (
	"fly/pkg/logging"
	"fmt"
	"gopkg.in/mgo.v2"
)

type Conf struct {
	Address     string // 地址 host:port
	MaxPoolSize int    // 连接池最大值
}

var (
	session *mgo.Session
)

// Init 初始化
func Init(c Conf) (err error) {
	if c.Address == "" {
		logging.Log.Warn("Init MongoDB not config")
		return
	}
	if c.MaxPoolSize <= 0 {
		c.MaxPoolSize = 10
	}
	session, err = mgo.Dial(fmt.Sprintf("mongodb://%s?maxPoolSize=%d", c.Address, c.MaxPoolSize))
	if err != nil {
		return
	}
	// 单调模式, 此模式下调用session必须使用copy/clone, 使用完进行close, 不然可能会有缓存冲突问题
	// 参考 https://cardinfolink.github.io/2017/05/17/mgo-session/
	//session.SetMode(mgo.Monotonic, true)
	// 最终一致性 直接使用 NewCollection 就可以, 但提交顺序可能存在不一致, 每次都是获取新的连接
	session.SetMode(mgo.Eventual, true)
	// set pool limit
	session.SetPoolLimit(c.MaxPoolSize)
	return nil
}

// NewCollection  实例一个连接
func NewCollection(database, table string) *mgo.Collection {
	return session.DB(database).C(table)
}
