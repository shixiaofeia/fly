package mongo

import (
	"gopkg.in/mgo.v2"
)

type Conf struct {
	Address string
}

var (
	session *mgo.Session
)

// Init 初始化
func Init(c Conf) (err error) {
	session, err = mgo.Dial(c.Address)
	if err != nil {
		return
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return nil
}

// NewCollection  实例一个连接
func NewCollection(database, table string) *mgo.Collection {
	return session.DB(database).C(table)
}
