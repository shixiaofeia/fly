package httpcode

import (
	"bytes"
	"encoding/json"
	constants "fly/internal/const"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type (
	RequestCall interface {
		Call() (respBytes []byte, err error)
	}
	RequestMethod struct {
		Method      string
		ContentType string
	}
	RequestParam struct {
		Method   RequestMethod
		Url      string
		Data     map[string]interface{}
		Response interface{}
		Header   map[string]string
	}
)

// 请求方式
var (
	RequestGet      = RequestMethod{Method: "GET"}
	RequestPostJson = RequestMethod{Method: "POST", ContentType: "application/json;charset=UTF-8"}
	RequestPostForm = RequestMethod{Method: "POST", ContentType: "application/x-www-form-urlencoded"}
)

// NewCall
func NewCall(method RequestMethod, url string, data map[string]interface{}, response interface{}) RequestCall {
	return NewCallWithHeader(method, url, data, response, nil)
}

// NewCallByGet
func NewCallByGet(url string) RequestCall {
	return NewCallWithHeader(RequestGet, url, nil, nil, nil)
}

// NewCallWithHeader
func NewCallWithHeader(method RequestMethod, url string, data map[string]interface{}, response interface{}, header map[string]string) RequestCall {
	return &RequestParam{Method: method, Url: url, Data: data, Response: response, Header: header}
}

// Call 调用接口
func (p *RequestParam) Call() (respBytes []byte, err error) {
	var (
		reader io.Reader
	)
	// 加载参数
	if p.Method == RequestGet {
		reader = bytes.NewReader([]byte{})
		if p.Data != nil {
			dataByte, _ := json.Marshal(p.Data)
			reader = bytes.NewReader(dataByte)
		}
	} else if p.Method == RequestPostForm {
		fromData := make(url.Values)
		for k, v := range p.Data {
			if val, ok := v.(string); ok {
				fromData[k] = []string{val}
			}
		}
		reader = strings.NewReader(fromData.Encode())
	}

	// 构建请求
	request, err := http.NewRequest(p.Method.Method, p.Url, reader)
	if err != nil {
		err = fmt.Errorf("NewRequest err: %v, param: %+v", err, p)
		return
	}
	// 防止Tcp复用
	request.Close = true
	if p.Method.ContentType != "" {
		request.Header.Set("Content-Type", p.Method.ContentType)
	}
	for k, v := range p.Header {
		request.Header.Add(k, v)
	}
	client := &http.Client{}
	client.Timeout = constants.CallInnerTimeOut
	resp, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("call api err: %v, param: %+v", err, p)
		return
	}
	defer request.Body.Close()

	defer resp.Body.Close()

	// 读取返回
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("ReadAll body err: %v, body: %s", err, resp.Body)
		return
	}
	// 需要结果则赋值结构体
	if p.Response != nil {
		_ = json.Unmarshal(respBytes, p.Response)
	}
	return
}
