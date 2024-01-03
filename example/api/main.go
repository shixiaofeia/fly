package main

import (
	"context"
	"fly/internal/config"
	"fly/internal/domain"
	"fly/pkg/clickhouse"
	"fly/pkg/logging"
	"fly/pkg/mysql"
	"fly/pkg/safego/safe"
	"fly/pkg/ws"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	err        error
	configPath = "configs/config.yml"
	wg         = new(sync.WaitGroup)
	gServer    = grpc.NewServer()
	srv        *http.Server
)

func main() {
	safe.Go(ws.PrintSocketLength)
	// 初始化路由
	app := gin.Default()
	Index(app)
	// 初始化业务表
	domain.Init(mysql.NewWriteDB(), clickhouse.NewSession())

	// 监听端口
	logging.Info("Start Web Server ")
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
	// 初始化配置
	config.Init(configPath)

	// 注册mysql
	if err = mysql.Init(config.Config.Mysql.Read, config.Config.Mysql.Write); err != nil {
		logging.Fatal("init mysql service err: " + err.Error())
	}
}

// shutdown 优雅的关闭
func shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Info("Shutdown start")
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
