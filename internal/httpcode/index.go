package httpcode

import (
	"bytes"
	"fly/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Req struct {
	ctx       *gin.Context
	body      []byte
	requestId string
	Log       logging.Encoder
}

// Ok 正确的json返回.
func (slf *Req) Ok(data interface{}) {
	slf.Code(Code200, nil, data)
}

// ParamError json返回参数错误.
func (slf *Req) ParamError(err error) {
	slf.Code(ParamErr, err, nil)
}

// ServiceError 通用错误处理.
func (slf *Req) ServiceError(err error) {
	slf.Code(ServiceErr, err, nil)
}

// Code 自定义code码返回.
func (slf *Req) Code(code ErrCode, err error, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	if err != nil {
		slf.Log.Warnf("api: %s, param: %s, err: %v", slf.ctx.Request.RequestURI, slf.body, err)
	}
	var runTime string
	if startTime, ok := slf.ctx.Get(CtxStartTime); ok {
		runTime = time.Now().Sub(startTime.(time.Time)).String()
	}

	slf.ctx.Header(CtxRequestId, slf.requestId)
	slf.ctx.JSON(http.StatusOK, map[string]interface{}{"code": code.Code, "message": code.Msg, "run": runTime, "data": data})
}

// NewRequest 解析post传参.
func NewRequest(ctx *gin.Context, params interface{}) (r *Req, err error) {
	uid := uuid.NewV4().String()
	r = &Req{
		ctx:       ctx,
		requestId: uid,
		Log:       logging.GetEncoder().WithField(CtxRequestId, uid),
	}
	if params != nil {
		defer func() {
			if err != nil {
				r.ParamError(err)
			}
		}()

		switch r.ctx.Request.Method {
		case http.MethodGet:
			err = ctx.ShouldBindQuery(params)
		case http.MethodDelete:
			err = ctx.ShouldBindQuery(params)
		case http.MethodPost:
			r.body, _ = ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(r.body))
			err = ctx.ShouldBindJSON(params)
		default:
		}
		if err != nil {
			return r, fmt.Errorf("api: %s getBody err: %v, body: %s", ctx.Request.RequestURI, err, string(r.body))
		}
	}
	return
}

// ToExcel 数据导出excel.
func (slf *Req) ToExcel(titleList []string, dataList interface{}, fileName string) {
	buf, _ := ExportExcel(titleList, dataList)
	content := bytes.NewReader(buf.Bytes())
	http.ServeContent(slf.ctx.Writer, slf.ctx.Request, fileName, time.Now(), content)
}

// ToSecondaryTitleExcel 导出二级标题.
func (slf *Req) ToSecondaryTitleExcel(firstTitle []string, secondTitle [][]string, dataList interface{}, fileName string) {
	buf, _ := ExportSecondaryTitleExcel(firstTitle, secondTitle, dataList)
	content := bytes.NewReader(buf.Bytes())
	http.ServeContent(slf.ctx.Writer, slf.ctx.Request, fileName, time.Now(), content)
}
