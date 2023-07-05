package sqldb

import (
	"gorm.io/gorm"
)

// InitTables 初始化创建MySQL数据库.
func InitTables(db *gorm.DB) {
	if db != nil {
		NewDemoSearch(db).CreateTable()
	}
}
