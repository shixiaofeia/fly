package example

import (
	"fly/example/api/example/controller"
	"fly/pkg/safego/safe"
	"github.com/gin-gonic/gin"
)

// InitApi 初始化路由.
func InitApi(app *gin.RouterGroup) {
	InitCtr()
	app.POST("/hello", controller.Hello)        // 请求示例
	app.GET("/export", controller.Export)       // 导出示例
	app.GET("/socket", controller.SocketHealth) // socket示例

	_demo := app.Group("/demo")
	{
		_demo.POST("/create", demoController.DemoCreate)   // 创建
		_demo.POST("/records", demoController.DemoRecords) // 列表
		_demo.POST("/info", demoController.DemoInfo)       // 详情
		_demo.POST("/update", demoController.DemoUpdate)   // 更新
		_demo.POST("/delete", demoController.DemoDelete)   // 删除
	}

	safe.Go(controller.GroupSend)
}
