package controller

import "fly/example/api/example/service"

var (
	demoService = service.NewDemoService()
)

type (
	DemoController struct{}
)

func NewDemoController() *DemoController {
	return &DemoController{}
}
