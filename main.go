package main

import (
	"fly/config"
	"fly/flydb/sqldb"
	"fly/interface"
	"fly/internal/monitor"
	"fly/pkg/logging"
	"fly/pkg/mq"
	"fly/pkg/mysql"
	"fly/pkg/redis"
	"github.com/kataras/iris/v12"
)

func main() {
	var err error
	// 初始化
	config.Init()
	// 注册mysql服务
	if err = mysql.Init(config.Config.Mysql.Read, config.Config.Mysql.Write); err != nil {
		logging.Log.Fatal("init mysql service err: " + err.Error())
	}
	if err = redis.Init(config.Config.Redis); err != nil {
		logging.Log.Fatal("init redis service err: " + err.Error())
	}
	//注册RabbitMQ服务
	err = mq.Init(config.Config.RabbitMq)
	if err != nil {
		logging.Log.Fatal("init rabbit mq err: " + err.Error())
	}
	// 初始化Mysql Tables
	sqldb.InitCreateTables()

	// 监控服务
	go monitor.InitMonitor()
	// 创建路由
	app := iris.New()
	api.Index(app)
	//监听端口
	logging.Log.Info("Start Fly Server API ")
	if err = app.Run(iris.Addr(":" + config.Config.ServerPort)); err != nil {
		logging.Log.Error("Start Fly Server API err: " + err.Error())
	}

}
