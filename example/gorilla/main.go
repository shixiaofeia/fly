package main

import (
	"fly/pkg/safego/safe"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	app := gin.Default()
	app.GET("/websocket_endpoint", Hello)
	log.Fatal(app.Run(":9999"))
}

func Hello(ctx *gin.Context) {
	upgrade := websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	safe.Go(func() {
		production(conn)
	})
	safe.Go(func() {
		consumer(conn)
	})
}

// production 生产.
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

// consumer 消费.
func consumer(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("consumer err: " + err.Error())
			return
		}
		fmt.Println(fmt.Sprintf("consumer data: %v", string(data)))
	}
}
