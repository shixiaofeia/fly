package controller

import (
	"fly/api/example/model"
	"fly/pkg/httpcode"
	"fmt"

	"github.com/kataras/iris/v12"
)

// Hello 请求示例.
func Hello(ctx iris.Context) {
	req := &model.HelloReq{}
	r, ok := httpcode.NewRequest(ctx, req)
	if !ok {
		return
	}
	res := model.HelloRes{
		Name: fmt.Sprintf("hello, %s", req.Name),
	}
	r.Log.Info(res.Name)
	r.JsonOk(res)
}

// Export 导出示例.
func Export(ctx iris.Context) {
	r, ok := httpcode.NewRequest(ctx, nil)
	if !ok {
		return
	}
	dataList := []*model.ExportRes{
		{1, "1", 2},
		{2, "3", 4},
	}
	r.ToExcel([]string{"Id", "年龄"}, dataList, "test")
}
