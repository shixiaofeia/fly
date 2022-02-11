package domain

import "fly/internal/domain/sqldb"

// InitDomain 初始化所有表.
func InitDomain() {
	sqldb.InitTables()
}
