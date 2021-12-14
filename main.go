package main

import (
	"context"
	"fly/config"
	"fly/domain"
	api "fly/interface"
	"fly/internal/monitor"
	"fly/pkg/logging"
	"fly/pkg/mq"
	"fly/pkg/mysql"
	"fly/pkg/redis"
	recover2 "fly/pkg/safego/recover"
	"github.com/kataras/iris/v12"
	"sync"
	"time"
)

var (
	err         error
	ctx, cancel = context.WithCancel(context.Background())
	wg          = new(sync.WaitGroup)
	app         = iris.New()
)

func main() {
	defer wg.Wait()

	// 初始化业务表
	domain.InitDomain()

	// 监控服务
	recover2.SafeGo(func() {
		monitor.InitMonitor(ctx)
	})

	// 初始化路由
	api.Index(app)

	// 监听端口
	logging.Log.Info("Start Fly Server API ")
	if err = app.Run(iris.Addr(":"+config.Config.ServerPort), iris.WithoutInterruptHandler); err != nil {
		logging.Log.Error("Start Fly Server API err: " + err.Error())
	}
}

func init() {
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
	if err = mq.Init(config.Config.RabbitMq); err != nil {
		logging.Log.Fatal("init rabbit mq err: " + err.Error())
	}

	// 优雅的关闭程序
	iris.RegisterOnInterrupt(func() {
		wg.Add(1)
		defer wg.Done()
		cancel()
		time.Sleep(5 * time.Second)
		// 关闭所有主机
		_ = app.Shutdown(ctx)
	})
}
