package main

import (
	"context"
	"flag"
	"fly/internal/api"
	"fly/internal/config"
	"fly/internal/domain"
	"fly/internal/monitor"
	"fly/internal/rpc"
	"fly/pkg/clickhouse"
	"fly/pkg/logging"
	"fly/pkg/mq"
	"fly/pkg/mysql"
	"fly/pkg/redis"
	"fly/pkg/safego/safe"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	err         error
	configPath  string
	ctx, cancel = context.WithCancel(context.Background())
	wg          = new(sync.WaitGroup)
	gServer     = grpc.NewServer()
	srv         *http.Server
)

func main() {
	// 初始化业务表
	domain.Init(mysql.NewWriteDB(), clickhouse.NewSession())

	// 监控服务
	safe.Go(func() {
		monitor.Start(ctx, wg)
	})

	// 初始化路由
	app := gin.Default()
	api.Index(app)

	// rpc
	initRpc()

	// 监听端口
	logging.Info("start Web Server")
	srv = &http.Server{
		Addr:    config.Config.Addr + ":" + config.Config.Port,
		Handler: app,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Fatal("start Web Server err: " + err.Error())
		}
	}()

	shutdown()
}

func init() {
	flag.StringVar(&configPath, "config", "./configs/config.yml", "配置文件路径以及文件名(必填)")
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
}

// initRpc 初始化rpc.
func initRpc() {
	rpc.Index(gServer)
	safe.Go(func() {
		lis, err := net.Listen("tcp", config.Config.Addr+":"+config.Config.RpcPort)
		if err != nil {
			logging.Fatal("start Rpc Listen err: " + err.Error())
		}
		logging.Info("start Rpc Server ")
		if err = gServer.Serve(lis); err != nil {
			logging.Fatal("start Rpc Server err: " + err.Error())
		}
	})
}

// shutdown 优雅的关闭
func shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Info("Shutdown start")
	cancel()
	// 关闭所有主机
	gServer.Stop()
	shutDownCtx, ShutDownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ShutDownCancel()
	if err = srv.Shutdown(shutDownCtx); err != nil {
		logging.Errorf("server shutdown err: %v", err)
	}
	wg.Wait()
	logging.Sync()
}
