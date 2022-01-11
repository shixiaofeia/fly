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
	"fly/rpc"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
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
	recover2.SafeGo(func() {
		monitor.InitMonitor(ctx)
	})

	// 初始化路由
	api.Index(app)

	// rpc
	initRpc()

	// 监听端口
	logging.Log.Info("Start Fly Server API ")
	if err = app.Run(iris.Addr(config.Config.ServerPort), iris.WithoutInterruptHandler); err != nil {
		logging.Log.Fatal("Start Fly Server API err: " + err.Error())
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

// initRpc 初始化rpc
func initRpc() {
	rpc.Index(gServer)
	recover2.SafeGo(func() {
		lis, err := net.Listen("tcp", config.Config.RpcPort)
		if err != nil {
			logging.Log.Fatal("Start Fly Rpc Listen err: " + err.Error())
		}
		logging.Log.Info("Start Fly Rpc Server ")
		if err = gServer.Serve(lis); err != nil {
			logging.Log.Fatal("Start Fly Rpc Server err: " + err.Error())
		}
	})
}
