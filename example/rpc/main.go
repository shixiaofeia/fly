package main

import (
	"fly/internal/config"
	"fly/pkg/logging"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var (
	configPath = "configs/config.json"
	wg         = new(sync.WaitGroup)
	gServer    = grpc.NewServer()
)

func main() {
	defer func() {
		wg.Wait()
		logging.Sync()
	}()

	Index(gServer)
	lis, err := net.Listen("tcp", ":"+config.Config.RpcPort)
	if err != nil {
		logging.Fatal("Start Rpc Listen err: " + err.Error())
	}
	logging.Info("Start Rpc Server ")
	if err = gServer.Serve(lis); err != nil {
		logging.Fatal("Start Rpc Server err: " + err.Error())
	}

}

func init() {
	// 初始化配置
	config.Init(configPath)
}
