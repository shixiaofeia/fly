package clickhouse

import (
	"fmt"

	"github.com/mailru/dbr"
	_ "github.com/mailru/go-clickhouse"
)

var (
	connect *dbr.Connection
)

type Config struct {
	Address  string
	Database string
}

// Init 初始化.
func Init(c Config) (err error) {
	connect, err = dbr.Open("clickhouse", fmt.Sprintf("%s/%s", c.Address, c.Database), nil)
	return
}

// NewSession  实例一个会话.
func NewSession() *dbr.Session {
	return connect.NewSession(nil)
}
