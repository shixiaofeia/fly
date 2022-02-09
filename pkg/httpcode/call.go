package httpcode

import (
	"bytes"
	"encoding/json"
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
		AddHeader(key, val string)
	}
	RequestParam struct {
		method      string
		contentType string
		url         string
		headers     map[string]string
	}

	RequestGet struct {
		RequestParam
	}
	RequestPostJson struct {
		RequestParam
		Body map[string]interface{}
	}
	RequestPostForm struct {
		RequestParam
		Body map[string]string
	}
)

func NewRequestGet(url string) RequestCall {
	return &RequestGet{RequestParam{method: "GET", url: url}}
}

func NewRequestPostJson(url string, body map[string]interface{}) RequestCall {
	return &RequestPostJson{RequestParam{method: "POST", contentType: "application/json;charset=UTF-8", url: url},
		body}
}

func NewRequestPostForm(url string, body map[string]string) RequestCall {
	return &RequestPostForm{RequestParam{method: "POST", contentType: "application/x-www-form-urlencoded", url: url},
		body}
}

// AddHeader 添加header头
func (slf *RequestParam) AddHeader(key, val string) {
	if slf.headers == nil {
		slf.headers = make(map[string]string)
	}
	slf.headers[key] = val
}

// Call
func (slf *RequestGet) Call() (respBytes []byte, err error) {
	// 构建请求
	request, err := http.NewRequest(slf.method, slf.url, nil)
	if err != nil {
		err = fmt.Errorf("NewRequest err: %v, param: %+v", err, slf)
		return
	}

	return call(request, slf.headers)
}

// Call
func (slf *RequestPostJson) Call() (respBytes []byte, err error) {
	var (
		reader io.Reader
	)

	if slf.Body != nil {
		dataByte, _ := json.Marshal(slf.Body)
		reader = bytes.NewReader(dataByte)
	}

	// 构建请求
	request, err := http.NewRequest(slf.method, slf.url, reader)
	if err != nil {
		err = fmt.Errorf("NewRequest err: %v, param: %+v", err, slf)
		return
	}

	request.Header.Set("Content-Type", slf.contentType)

	return call(request, slf.headers)
}

// Call
func (slf *RequestPostForm) Call() (respBytes []byte, err error) {
	var (
		reader io.Reader
	)

	if slf.Body != nil {
		fromData := make(url.Values)
		for k, v := range slf.Body {
			fromData[k] = []string{v}
		}
		reader = strings.NewReader(fromData.Encode())
	}

	// 构建请求
	request, err := http.NewRequest(slf.method, slf.url, reader)
	if err != nil {
		err = fmt.Errorf("NewRequest err: %v, param: %+v", err, slf)
		return
	}

	request.Header.Set("Content-Type", slf.contentType)

	return call(request, slf.headers)
}

// call 发送请求
func call(request *http.Request, headers map[string]string) (respBytes []byte, err error) {
	var (
		resp *http.Response
	)

	// 防止Tcp复用
	request.Close = true
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	client := &http.Client{}
	client.Timeout = CallInnerTimeOut
	resp, err = client.Do(request)
	if err != nil {
		return
	}

	if request.Body != nil {
		defer request.Body.Close()
	}
	defer resp.Body.Close()

	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("ReadAll body err: %v, body: %s", err, resp.Body)
		return
	}
	return
}
