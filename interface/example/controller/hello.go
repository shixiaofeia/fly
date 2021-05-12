package controller

import (
	"fly/pkg/httpcode"
	"fmt"
	"github.com/kataras/iris/v12"
)

// Hello 请求示例
func Hello(ctx iris.Context) {
	type Param struct {
		Name string `json:"name" validate:"nonzero"`
	}
	req := &Param{}
	r, ok := httpcode.NewRequest(ctx, req)
	if !ok {
		return
	}
	type Response struct {
		Name string `json:"name"`
	}
	res := Response{
		Name: fmt.Sprintf("hello, %s", req.Name),
	}
	r.Log.Info(res.Name)
	r.JsonOk(res)
}

// Export 导出示例
func Export(ctx iris.Context) {
	r, ok := httpcode.NewRequest(ctx, nil)
	if !ok {
		return
	}
	type Response struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	dataList := []interface{}{
		Response{"1", 2},
		Response{"3", 4},
	}
	r.DataToExcel([]string{"名字", "年龄"}, dataList, "test")
}
