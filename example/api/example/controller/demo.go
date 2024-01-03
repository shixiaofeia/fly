package controller

import (
	"fly/example/api/example/model"
	"fly/internal/httpcode"
	"github.com/gin-gonic/gin"
)

// DemoCreate demo创建.
func (slf *DemoController) DemoCreate(ctx *gin.Context) {
	var (
		req    = new(model.DemoCreateReq)
		r, err = httpcode.NewRequest(ctx, req)
	)

	if err != nil {
		return
	}

	if err = slf.srv.DemoCreate(req); err != nil {
		r.ServiceError(err)
		return
	}

	r.Ok(nil)
}

// DemoRecords demo列表.
func (slf *DemoController) DemoRecords(ctx *gin.Context) {
	var (
		req    = new(model.DemoRecordReq)
		r, err = httpcode.NewRequest(ctx, req)
	)

	if err != nil {
		return
	}

	res, err := slf.srv.DemoRecords(req)
	if err != nil {
		r.ServiceError(err)
		return
	}

	r.Ok(res)
}

// DemoInfo demo详情.
func (slf *DemoController) DemoInfo(ctx *gin.Context) {
	var (
		req    = new(model.DemoInfoReq)
		r, err = httpcode.NewRequest(ctx, req)
	)

	if err != nil {
		return
	}

	res, err := slf.srv.DemoInfo(req)
	if err != nil {
		r.ServiceError(err)
		return
	}

	r.Ok(res)
}

// DemoUpdate demo更新.
func (slf *DemoController) DemoUpdate(ctx *gin.Context) {
	var (
		req    = new(model.DemoUpdateReq)
		r, err = httpcode.NewRequest(ctx, req)
	)

	if err != nil {
		return
	}

	if code, err := slf.srv.DemoUpdate(req); err != nil {
		r.Code(code, err, nil)
		return
	}

	r.Ok(nil)
}

// DemoDelete demo删除.
func (slf *DemoController) DemoDelete(ctx *gin.Context) {
	var (
		req    = new(model.DemoDeleteReq)
		r, err = httpcode.NewRequest(ctx, req)
	)

	if err != nil {
		return
	}

	if err = slf.srv.DemoDelete(req); err != nil {
		r.ServiceError(err)
		return
	}

	r.Ok(nil)
}
