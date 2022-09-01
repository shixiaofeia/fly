package main

import (
	"context"
	"flag"
	"fly/internal/api"
	"fly/internal/config"
	"fly/internal/domain"
	"fly/internal/monitor"
	"fly/internal/rpc"
	"fly/pkg/logging"
	"fly/pkg/mq"
	"fly/pkg/mysql"
	"fly/pkg/redis"
	"fly/pkg/safego/safe"
	"net"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
)

var (
	err         error
	configPath  string
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

	// 初始化业务表
	domain.Init()

	// 监控服务
	safe.Go(func() {
		monitor.Start(ctx)
	})

	// 初始化路由
	api.Index(app)

	// rpc
	initRpc()

	// 监听端口
	logging.Info("start Web Server")
	if err = app.Run(iris.Addr(":"+config.Config.Port), iris.WithoutInterruptHandler); err != nil {
		logging.Fatal("start Web Server err: " + err.Error())
	}
}

func init() {
	flag.StringVar(&configPath, "config", "./configs/config.json", "配置文件路径以及文件名(必填)")
	flag.Parse()

	// 初始化配置
	config.Init(configPath)

	// 注册mysql
	if err = mysql.Init(config.Config.Mysql.Read, config.Config.Mysql.Write); err != nil {
		logging.Fatal("init mysql service err: " + err.Error())
	}

	// 注册redis
	if err = redis.Init(config.Config.Redis); err != nil {
		logging.Fatal("init redis service err: " + err.Error())
	}

	// 注册RabbitMQ
	if err = mq.Init(config.Config.Mq); err != nil {
		logging.Fatal("init rabbit mq err: " + err.Error())
	}

	// 优雅关闭程序
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

// initRpc 初始化rpc.
func initRpc() {
	rpc.Index(gServer)
	safe.Go(func() {
		lis, err := net.Listen("tcp", ":"+config.Config.RpcPort)
		if err != nil {
			logging.Fatal("start Rpc Listen err: " + err.Error())
		}
		logging.Info("start Rpc Server ")
		if err = gServer.Serve(lis); err != nil {
			logging.Fatal("start Rpc Server err: " + err.Error())
		}
	})
}
