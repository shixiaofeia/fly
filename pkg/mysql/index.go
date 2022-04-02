package mysql

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	writeDB     *gorm.DB
	readDB      *gorm.DB
	mysqlConfig = mysql.Config{
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	newLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // Log level
			Colorful:                  true,          // 彩色打印
			IgnoreRecordNotFoundError: true,          // 忽略记录为空的错误
		},
	)
)

// Init 初始化db.
func Init(readConf, writeConf Conf) (err error) {
	if err = InitReadDB(readConf); err != nil {
		return
	}
	if err = InitWriteDB(writeConf); err != nil {
		return
	}
	return
}

// NewReadDB 只读.
func NewReadDB() *gorm.DB {
	return readDB
}

// NewWriteDB 写.
func NewWriteDB() *gorm.DB {
	return writeDB
}

// InitReadDB 初始化读.
func InitReadDB(c Conf) (err error) {
	if c.Address == "" {
		return
	}
	dsn := c.User + ":" + c.Password + "@tcp(" + c.Address + ":" + c.Port + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig.DSN = dsn
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{Logger: newLogger})
	if err != nil {
		return fmt.Errorf("InitReadDB err: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdleCoon)
	sqlDB.SetMaxOpenConns(c.MaxOpenCoon)
	readDB = db
	return
}

// InitWriteDB 初始化写.
func InitWriteDB(c Conf) (err error) {
	if c.Address == "" {
		return
	}
	dsn := c.User + ":" + c.Password + "@tcp(" + c.Address + ":" + c.Port + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig.DSN = dsn
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{Logger: newLogger})
	if err != nil {
		return fmt.Errorf("InitWriteDB err: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdleCoon)
	sqlDB.SetMaxOpenConns(c.MaxOpenCoon)
	// 回调函数
	_ = db.Callback().Create().Before("gorm:create").Register("beforeCreateUpTime", beforeCreateUpTime)
	_ = db.Callback().Update().Before("gorm:update").Register("beforeUpdateUpTime", beforeUpdateUpTime)
	writeDB = db
	return
}

// beforeCreateUpTime 在插入之前更新时间戳.
func beforeCreateUpTime(tx *gorm.DB) {
	tx.Statement.SetColumn("create_time", time.Now().Unix())
	tx.Statement.SetColumn("update_time", time.Now().Unix())
}

// beforeUpdateUpTime 在更新之前更新时间戳.
func beforeUpdateUpTime(tx *gorm.DB) {
	tx.Statement.SetColumn("update_time", time.Now().Unix())
}
