package domain

import "fly/internal/domain/sqldb"

// Init 初始化所有表.
func Init() {
	sqldb.InitTables()
}
