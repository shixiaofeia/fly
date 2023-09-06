package mongo

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

type (
	Conf struct {
		Host        string
		Port        string
		DataBase    string
		MaxPoolSize int // 连接池最大值
	}
	GridFS struct {
		client *mgo.GridFS
	}
)

var (
	_database *mgo.Database
)

// Init 初始化.
func Init(c Conf) error {
	if c.Host == "" || c.Port == "" {
		return nil
	}
	if c.MaxPoolSize <= 0 {
		c.MaxPoolSize = 10
	}
	session, err := mgo.Dial(fmt.Sprintf("mongodb://%s:%s?maxPoolSize=%d", c.Host, c.Port, c.MaxPoolSize))
	if err != nil {
		return err
	}
	if err = session.Ping(); err != nil {
		return err
	}

	// 单调模式, 此模式下调用session必须使用copy/clone, 使用完进行close, 不然可能会有缓存冲突问题
	// 参考 https://cardinfolink.github.io/2017/05/17/mgo-session/
	// session.SetMode(mgo.Monotonic, true)
	// 最终一致性 直接使用 NewCollection 就可以, 但提交顺序可能存在不一致, 每次都是获取新的连接
	session.SetMode(mgo.Eventual, true)
	// set pool limit
	session.SetPoolLimit(c.MaxPoolSize)
	_database = session.DB(c.DataBase)

	return nil
}

// NewCollection  实例一个连接.
func NewCollection(table string) *mgo.Collection {
	return _database.C(table)
}

// NewGridFS 实例GridFS, 主要用来存取大文件, 比如图片/视频/音频等.
func NewGridFS(prefix string) *GridFS {
	return &GridFS{client: _database.GridFS(prefix)}
}
