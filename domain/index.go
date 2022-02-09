package domain

import "fly/domain/sqldb"

// InitDomain 初始化所有表.
func InitDomain() {
	sqldb.InitTables()
}
