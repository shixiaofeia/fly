package httpcode

import (
	"bytes"
	"encoding/json"
	"fly/internal/const"
	"fly/pkg/logging"
	"fmt"
	"github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
	"gopkg.in/validator.v2"
	"reflect"
	"strings"
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
	Log       *logging.PrefixLog
}

// JsonOk 正确的json返回
func (r *Req) JsonOk(data interface{}) {
	r.JsonCode(Code200, data)
}

// JsonParamError json返回参数错误
func (r *Req) JsonParamError() {
	r.JsonCode(ParamErr, nil)
}

// JsonMysqlError 数据库处理错误
func (r *Req) JsonMysqlError() {
	r.JsonCode(MysqlErr, nil)
}

// JsonCode 自定义code码返回
func (r *Req) JsonCode(code ErrCode, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	startTime := r.ctx.Values().GetInt64Default(CtxStartTime, time.Now().UnixNano())
	takeTime := (time.Now().UnixNano() - startTime) / 1e6
	r.ctx.Header(CtxRequestId, r.requestId)
	_, _ = r.ctx.JSON(map[string]interface{}{"code": code.Code, "message": code.Msg, "run": takeTime, "data": data})
	r.Log.Info(fmt.Sprintf("api: %s, param: %s, code: %+v, response: %+v", r.ctx.Request().RequestURI, r.body, code, data))
}

// NewRequest 解析post传参
func NewRequest(ctx iris.Context, params interface{}) (r *Req, b bool) {
	uid := uuid.NewV4().String()
	r = &Req{
		ctx:       ctx,
		requestId: uid,
		Log:       logging.NewPrefixLog("API X-Request-Id: " + uid),
	}
	if params != nil {
		body, err := ctx.GetBody()
		if err != nil {
			r.Log.Error(fmt.Sprintf("api: %s GetBody err: %v", ctx.Request().RequestURI, err))
			r.JsonParamError()
			return
		}
		if len(body) > 0 {
			if err = json.Unmarshal(body, params); err != nil {
				r.Log.Error(fmt.Sprintf("api: %s Unmarshal err: %v, body: %s", ctx.Request().RequestURI, err, body))
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

// DataToExcel 数据导出excel, dataList里面的对象为指针
func (r *Req) DataToExcel(titleList []string, data interface{}, fileName string) {
	// 生成一个新的文件
	file := xlsx.NewFile()
	// 添加sheet页
	sheet, _ := file.AddSheet("Sheet1")
	// 插入表头
	titleRow := sheet.AddRow()
	for _, v := range titleList {
		cell := titleRow.AddCell()
		cell.Value = v
		cell.GetStyle().Font.Color = "00FF0000"
	}
	// 插入内容
	for _, v := range data.([]interface{}) {
		row := sheet.AddRow()
		dataList := make([]interface{}, 0)
		RecursionStructValueToSlice(v, reflect.Value{}, &dataList)
		row.WriteSlice(&dataList, -1)
	}
	fileName = fmt.Sprintf("%s.xlsx", fileName)
	// 打开save即可本地保存文件
	_ = file.Save(fileName)
	r.ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	r.ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	var buffer bytes.Buffer
	_ = file.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())
	_ = r.ctx.ServeContent(content, fileName, time.Now(), true)
}

// RecursionStructValueToSlice 递归多层嵌套结构体的值转入slice
func RecursionStructValueToSlice(val interface{}, v reflect.Value, data *[]interface{}) {
	// 判断是否为指针
	if val == nil {
	} else if reflect.ValueOf(val).Type().Kind() == reflect.Struct {
		v = reflect.ValueOf(val)
	} else {
		v = reflect.Indirect(reflect.ValueOf(val))
	}
	for i := 0; i < v.NumField(); i++ {
		// 判断是否是嵌套结构
		if v.Field(i).Type().Kind() == reflect.Struct {
			RecursionStructValueToSlice(nil, v.Field(i), data)
		} else {
			t := v.Type().Field(i)
			tags := ParseTagSetting(t.Tag, "excel")
			// 忽略不导出字段
			if _, ok := tags["-"]; ok {
				continue
			}
			*data = append(*data, v.Field(i).Interface())
		}
	}
	return
}

// ParseTagSetting 获取字段tags
func ParseTagSetting(tags reflect.StructTag, key ...string) map[string]string {
	setting := map[string]string{}
	for _, v := range key {
		str := tags.Get(v)
		if len(str) == 0 {
			continue
		}
		tagList := strings.Split(str, ";")
		for _, value := range tagList {
			if len(value) == 0 {
				continue
			}
			tagV := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(tagV[0]))
			if len(tagV) >= 2 {
				setting[k] = strings.Join(tagV[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}
