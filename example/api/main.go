package main

import (
	"context"
	"fly/internal/config"
	"fly/pkg/logging"
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
	defer wg.Wait()

	// 初始化路由
	Index(app)

	// 监听端口
	logging.Log.Info("Start Web Server ")
	if err = app.Run(iris.Addr(config.Config.ServerPort), iris.WithoutInterruptHandler); err != nil {
		logging.Log.Fatal("Start Web Server err: " + err.Error())
	}
}

func init() {
	// 初始化配置
	config.Init(configPath)

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
