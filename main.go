package main

import (
	"context"
	"fly/api"
	"fly/config"
	"fly/domain"
	"fly/internal/monitor"
	"fly/pkg/logging"
	"fly/pkg/mq"
	"fly/pkg/mysql"
	"fly/pkg/redis"
	"fly/pkg/safego/safe"
	"fly/rpc"
	"net"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
)

var (
	err         error
	ctx, cancel = context.WithCancel(context.Background())
	wg          = new(sync.WaitGroup)
	app         = iris.New()
	gServer     = grpc.NewServer()
)

func main() {
	defer wg.Wait()

	// 初始化业务表
	domain.InitDomain()

	// 监控服务
	safe.Go(func() {
		monitor.InitMonitor(ctx)
	})

	// 初始化路由
	api.Index(app)

	// rpc
	initRpc()

	// 监听端口
	logging.Log.Info("Start Web Server ")
	if err = app.Run(iris.Addr(config.Config.ServerPort), iris.WithoutInterruptHandler); err != nil {
		logging.Log.Fatal("Start Web Server err: " + err.Error())
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
		gServer.Stop()
		_ = app.Shutdown(ctx)
	})
}

// initRpc 初始化rpc.
func initRpc() {
	rpc.Index(gServer)
	safe.Go(func() {
		lis, err := net.Listen("tcp", config.Config.RpcPort)
		if err != nil {
			logging.Log.Fatal("Start Rpc Listen err: " + err.Error())
		}
		logging.Log.Info("Start Rpc Server ")
		if err = gServer.Serve(lis); err != nil {
			logging.Log.Fatal("Start Rpc Server err: " + err.Error())
		}
	})
}
