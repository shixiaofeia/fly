package domain

import (
	"fly/internal/domain/sqldb"
	"gorm.io/gorm"
)

// Init 初始化所有表.
func Init(orm *gorm.DB) {
	sqldb.InitTables(orm)
}
