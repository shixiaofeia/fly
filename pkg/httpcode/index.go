package httpcode

import (
	"bytes"
	"encoding/json"
	"fly/internal/const"
	"fly/pkg/logging"
	"github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gopkg.in/validator.v2"
	"reflect"
	"time"
)

const (
	CtxStartTime = "startTime"    // 运行起始时间
	CtxRequestId = "X-Request-Id" // 请求唯一id
)

type Req struct {
	ctx       iris.Context
	body      []byte
	requestId string
	Log       *zap.SugaredLogger
}

// JsonOk 正确的json返回
func (slf *Req) JsonOk(data interface{}) {
	slf.JsonCode(Code200, data)
}

// JsonParamError json返回参数错误
func (slf *Req) JsonParamError() {
	slf.JsonCode(ParamErr, nil)
}

// JsonServiceError 通用错误处理
func (slf *Req) JsonServiceError() {
	slf.JsonCode(ServiceErr, nil)
}

// JsonCode 自定义code码返回
func (slf *Req) JsonCode(code ErrCode, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	var runTime string
	if startTime, ok := slf.ctx.Values().Get(CtxStartTime).(time.Time); ok {
		runTime = time.Now().Sub(startTime).String()
	}

	slf.ctx.Header(CtxRequestId, slf.requestId)
	_, _ = slf.ctx.JSON(map[string]interface{}{"code": code.Code, "message": code.Msg, "run": runTime, "data": data})
	slf.Log.Infof("api: %s, run: %s, param: %s, code: %d", slf.ctx.Request().RequestURI, runTime, slf.body, code.Code)
}

// NewRequest 解析post传参
func NewRequest(ctx iris.Context, params interface{}) (r *Req, b bool) {
	uid := uuid.NewV4().String()
	r = &Req{
		ctx:       ctx,
		requestId: uid,
		Log:       logging.NewWithField("X-Request-Id", uid),
	}
	if params != nil {
		body, err := ctx.GetBody()
		if err != nil {
			r.Log.Errorf("api: %s GetBody err: %v", ctx.Request().RequestURI, err)
			r.JsonParamError()
			return
		}
		if len(body) > 0 {
			if err = json.Unmarshal(body, params); err != nil {
				r.Log.Errorf("api: %s Unmarshal err: %v, body: %s", ctx.Request().RequestURI, err, body)
				r.JsonParamError()
				return
			}
			r.body = body
		}
		// 页数限制
		if page := reflect.Indirect(reflect.ValueOf(params)).FieldByName("Page"); page.String() != "<invalid Value>" && page.Int() > constants.MaxPage {
			r.JsonCode(PageErr, nil)
			return
		}
		// 条数限制
		if size := reflect.Indirect(reflect.ValueOf(params)).FieldByName("Size"); size.String() != "<invalid Value>" && size.Int() > constants.MaxSize {
			r.JsonCode(SizeErr, nil)
			return
		}
		// 参数校验
		if err = validator.Validate(params); err != nil {
			r.Log.Error("Validate param err: " + err.Error())
			r.JsonParamError()
			return
		}
	}
	b = true
	return
}

// ToExcel 数据导出excel
func (slf *Req) ToExcel(titleList []string, dataList interface{}, fileName string) {
	buf, _ := ExportExcel(titleList, dataList)
	content := bytes.NewReader(buf.Bytes())
	_ = slf.ctx.ServeContent(content, fileName, time.Now(), true)
}

// ToSecondaryTitleExcel 导出二级标题
func (slf *Req) ToSecondaryTitleExcel(firstTitle []string, secondTitle [][]string, dataList interface{}, fileName string) {
	buf, _ := ExportSecondaryTitleExcel(firstTitle, secondTitle, dataList)
	content := bytes.NewReader(buf.Bytes())
	_ = slf.ctx.ServeContent(content, fileName, time.Now(), true)
}
