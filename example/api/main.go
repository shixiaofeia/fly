package main

import (
	"context"
	"fly/internal/config"
	"fly/internal/domain"
	"fly/pkg/logging"
	"fly/pkg/mysql"
	"fly/pkg/safego/safe"
	"fly/pkg/ws"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
)

var (
	err         error
	configPath  = "configs/config.json"
	ctx, cancel = context.WithCancel(context.Background())
	wg          = new(sync.WaitGroup)
	app         = iris.New()
	gServer     = grpc.NewServer()
)

func main() {
	defer func() {
		wg.Wait()
		logging.Sync()
	}()
	safe.Go(ws.PrintSocketLength)
	// 初始化路由
	Index(app)
	// 初始化业务表
	domain.Init()

	// 监听端口
	logging.Info("Start Web Server ")
	if err = app.Run(iris.Addr(":"+config.Config.Port), iris.WithoutInterruptHandler); err != nil {
		logging.Fatal("Start Web Server err: " + err.Error())
	}
}

func init() {
	// 初始化配置
	config.Init(configPath)

	// 注册mysql
	if err = mysql.Init(config.Config.Mysql.Read, config.Config.Mysql.Write); err != nil {
		logging.Fatal("init mysql service err: " + err.Error())
	}

	// 优雅的关闭程序
	iris.RegisterOnInterrupt(func() {
		wg.Add(1)
		defer wg.Done()
		cancel()
		time.Sleep(5 * time.Second)
		// 关闭所有主机
		gServer.Stop()
		_ = app.Shutdown(ctx)
	})
}
