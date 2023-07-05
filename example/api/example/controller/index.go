package controller

import "fly/example/api/example/service"

type (
	DemoController struct {
		srv *service.DemoService
	}
)

func NewDemoController(srv *service.DemoService) *DemoController {
	return &DemoController{srv: srv}
}
