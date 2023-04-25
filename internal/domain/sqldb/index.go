package sqldb

import "fly/pkg/mysql"

// InitTables 初始化创建MySQL数据库.
func InitTables() {
	var db = mysql.NewWriteDB()
	if db != nil {
		NewDemoSearch(db).CreateTable()
	}
}
