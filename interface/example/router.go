package example

import (
	"fly/interface/example/controller"
	"github.com/kataras/iris/v12"
)

// InitApi
func InitApi(app iris.Party) {
	app.Post("/hello", controller.Hello)        // 请求示例
	app.Get("/export", controller.Export)       // 导出示例
	app.Get("/socket", controller.SocketHealth) // socket示例
	go controller.GroupSend()
}
