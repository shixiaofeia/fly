package domain

import (
	"fly/internal/domain/ckdb"
	"fly/internal/domain/sqldb"
	"github.com/mailru/dbr"
	"gorm.io/gorm"
)

// Init 初始化所有表.
func Init(orm *gorm.DB, session *dbr.Session) {
	sqldb.InitTables(orm)
	ckdb.InitTables(session)
}
