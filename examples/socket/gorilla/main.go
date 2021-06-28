package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
	"log"
	"net/http"
	"time"
)

func main() {
	app := iris.New()
	app.Get("/websocket_endpoint", Hello)
	log.Fatal(app.Run(iris.Addr(":9999")))
}

func Hello(ctx iris.Context) {
	upgrade := websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrade.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
	if err != nil {
		return
	}
	go production(conn)
	go consumer(conn)
}

// production 生产
func production(conn *websocket.Conn) {
	for {
		time.Sleep(1 * time.Second)

		if err := conn.WriteMessage(websocket.TextMessage, []byte("111")); err != nil {
			if err == websocket.ErrCloseSent {
				return
			}
			fmt.Println("err: " + err.Error())
		}
	}
}

// consumer 消费
func consumer(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if err.Error() == fmt.Sprintf("websocket: close %d (no status)", websocket.CloseNoStatusReceived) {
				return
			}
			fmt.Println("consumer err: " + err.Error())
			continue
		}
		fmt.Println(fmt.Sprintf("consumer data: %v", string(data)))
	}
}
