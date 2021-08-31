package main

import (
	"fmt"
	websocket2 "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kataras/neffos/gorilla"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

var ws *neffos.Server

func main() {
	upgrade := websocket2.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws = neffos.New(gorilla.Upgrader(upgrade), neffos.Namespaces{})
	app := iris.New()
	app.Get("/websocket_endpoint", Hello)
	log.Fatal(app.Run(iris.Addr(":9999")))
}

// Hello
func Hello(ctx iris.Context) {
	conn := websocket.Upgrade(ctx, func(ctx context.Context) string {
		return uuid.NewV4().String()
	}, ws)
	go production(conn)
	go consumer(conn)
}

// production 生产
func production(conn *neffos.Conn) {
	for {
		time.Sleep(1 * time.Second)
		if err := conn.Socket().WriteText([]byte("111"), 0); err != nil {
			if conn.IsClosed() {
				return
			}
			fmt.Println("err: " + err.Error())
		}
	}
}

// consumer 消费 TODO 存在漏接消息的情况, 已提[issues](https://github.com/kataras/neffos/issues/58)
func consumer(conn *neffos.Conn) {
	for {
		data, _, err := conn.Socket().ReadData(0)
		if err != nil {
			if conn.IsClosed() {
				return
			}
			fmt.Println("consumer err: " + err.Error())
			continue
		}
		fmt.Println(fmt.Sprintf("consumer data: %v", string(data)))
	}
}
