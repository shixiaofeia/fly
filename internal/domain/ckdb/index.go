package ckdb

import (
	"github.com/mailru/dbr"
)

// InitTables 初始化创建CK数据库.
func InitTables(session *dbr.Session) {
	if session != nil {
		_ = NewDemoSearch(session).CreateTable()
	}
}
