package example

import (
	"fly/example/api/example/controller"
	"fly/example/api/example/service"
	"fly/pkg/mysql"
)

var (
	demoController *controller.DemoController
)

func InitCtr() {
	var (
		orm = mysql.NewWriteDB()

		demoSrv = service.NewDemoService(orm)
	)

	demoController = controller.NewDemoController(demoSrv)
}
