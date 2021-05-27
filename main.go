package main

import (
	"context"
	"fly/config"
	"fly/flydb/sqldb"
	api "fly/interface"
	"fly/internal/monitor"
	"fly/pkg/logging"
	"fly/pkg/mq"
	"fly/pkg/mysql"
	"fly/pkg/redis"
	"github.com/kataras/iris/v12"
	"sync"
	"time"
)

func main() {
	var err error
	app := iris.New()

	// 优雅的关闭程序
	wg := new(sync.WaitGroup)
	defer wg.Wait()
	iris.RegisterOnInterrupt(func() {
		wg.Add(1)
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		// 关闭所有主机
		_ = app.Shutdown(ctx)
	})

	// 初始化配置
	config.Init()
	// 注册mysql
	if err = mysql.Init(config.Config.Mysql.Read, config.Config.Mysql.Write); err != nil {
		logging.Log.Fatal("init mysql service err: " + err.Error())
	}
	// 注册redis
	if err = redis.Init(config.Config.Redis); err != nil {
		logging.Log.Fatal("init redis service err: " + err.Error())
	}
	// 注册RabbitMQ
	err = mq.Init(config.Config.RabbitMq)
	if err != nil {
		logging.Log.Fatal("init rabbit mq err: " + err.Error())
	}
	// 初始化业务表
	sqldb.InitCreateTables()

	// 监控服务
	go monitor.InitMonitor()
	// 初始化路由
	api.Index(app)
	// 监听端口
	logging.Log.Info("Start Fly Server API ")
	if err = app.Run(iris.Addr(":"+config.Config.ServerPort), iris.WithoutInterruptHandler); err != nil {
		logging.Log.Error("Start Fly Server API err: " + err.Error())
	}

}
