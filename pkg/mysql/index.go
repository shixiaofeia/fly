package mysql

import (
	"context"
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
	createBatchSize = 200 // 批量插入条数
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
	if c.Addr == "" {
		return
	}
	dsn := c.User + ":" + c.Pwd + "@tcp(" + c.Addr + ":" + c.Port + ")/" + c.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig.DSN = dsn
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{Logger: newLogger})
	if err != nil {
		return fmt.Errorf("InitReadDB err: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdleCoon)
	sqlDB.SetMaxOpenConns(c.MaxOpenCoon)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = sqlDB.PingContext(ctx); err != nil {
		return err
	}
	readDB = db

	return
}

// InitWriteDB 初始化写.
func InitWriteDB(c Conf) (err error) {
	if c.Addr == "" {
		return
	}
	dsn := c.User + ":" + c.Pwd + "@tcp(" + c.Addr + ":" + c.Port + ")/" + c.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig.DSN = dsn
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{Logger: newLogger, CreateBatchSize: createBatchSize})
	if err != nil {
		return fmt.Errorf("InitWriteDB err: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdleCoon)
	sqlDB.SetMaxOpenConns(c.MaxOpenCoon)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = sqlDB.PingContext(ctx); err != nil {
		return err
	}
	writeDB = db

	return
}
