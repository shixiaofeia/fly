package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"log"
)

// [...]

func main() {
	onChat := func(ns *neffos.NSConn, msg neffos.Message) error {
		//ctx := websocket.GetContext(ns.Conn)
		fmt.Println(string(msg.Body))
		return nil
	}

	app := iris.New()
	ws := neffos.New(websocket.DefaultGorillaUpgrader, neffos.Namespaces{
		"default": neffos.Events{
			"chat": onChat,
		},
	})

	app.Get("/websocket_endpoint", websocket.Handler(ws))
	log.Fatal(app.Run(iris.Addr(":9999")))
}
