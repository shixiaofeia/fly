package example

import (
	"fly/example/api/example/controller"
	"fly/pkg/safego/safe"

	"github.com/kataras/iris/v12"
)

// InitApi 初始化路由.
func InitApi(app iris.Party) {
	InitCtr()
	app.Post("/hello", controller.Hello)        // 请求示例
	app.Get("/export", controller.Export)       // 导出示例
	app.Get("/socket", controller.SocketHealth) // socket示例

	_demo := app.Party("/demo")
	{
		_demo.Post("/create", demoController.DemoCreate)   // 创建
		_demo.Post("/records", demoController.DemoRecords) // 列表
		_demo.Post("/info", demoController.DemoInfo)       // 详情
		_demo.Post("/update", demoController.DemoUpdate)   // 更新
		_demo.Post("/delete", demoController.DemoDelete)   // 删除
	}

	safe.Go(controller.GroupSend)
}
