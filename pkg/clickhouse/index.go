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
	Host     string
	Port     string
	Database string
	User     string
	Pwd      string
}

// Init 初始化.
func Init(c Config) (err error) {
	if c.Host == "" {
		return
	}

	dsn := fmt.Sprintf("http://%s:%s@%s:%s/%s", c.User, c.Pwd, c.Host, c.Port, c.Database)
	connect, err = dbr.Open("clickhouse", dsn, nil)
	return
}

// NewSession  实例一个会话.
func NewSession() *dbr.Session {
	if connect != nil {
		return connect.NewSession(nil)
	}
	return nil
}
