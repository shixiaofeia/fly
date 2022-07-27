package httpcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	RequestCall interface {
		AddHeader(key, val string)
		AddBasic(user, pwd string)
		Call() (respBytes []byte, err error)
	}
	BasicAuth struct {
		User string
		Pwd  string
	}
	RequestParam struct {
		method      string
		contentType string
		url         string
		headers     map[string]string
		basic       *BasicAuth
	}

	RequestGet struct {
		RequestParam
	}
	RequestPostJson struct {
		RequestParam
		Body interface{}
	}
	RequestPostForm struct {
		RequestParam
		Body map[string]string
	}
	RequestPostFormWithFile struct {
		RequestPostForm
		Files map[string][]byte
	}
)

func NewRequestGet(url string) RequestCall {
	return &RequestGet{RequestParam{method: http.MethodGet, url: url}}
}

func NewRequestPostJson(url string, body interface{}) RequestCall {
	return &RequestPostJson{RequestParam{method: http.MethodPost, contentType: "application/json;charset=UTF-8", url: url},
		body}
}

func NewRequestPostForm(url string, body map[string]string) RequestCall {
	return &RequestPostForm{RequestParam{method: http.MethodPost, contentType: "application/x-www-form-urlencoded", url: url},
		body}
}

func NewRequestPostFormWithFile(url string, body map[string]string, files map[string][]byte) RequestCall {
	return &RequestPostFormWithFile{RequestPostForm{
		RequestParam: RequestParam{method: http.MethodPost, url: url},
		Body:         body,
	}, files}
}

// AddHeader 添加header头.
func (slf *RequestParam) AddHeader(key, val string) {
	if slf.headers == nil {
		slf.headers = make(map[string]string)
	}
	slf.headers[key] = val

	return
}

func (slf *RequestParam) AddBasic(user, pwd string) {
	slf.basic = &BasicAuth{User: user, Pwd: pwd}

	return
}

// Call 调用.
func (slf *RequestGet) Call() (respBytes []byte, err error) {
	// 构建请求
	request, err := http.NewRequest(slf.method, slf.url, nil)
	if err != nil {
		err = fmt.Errorf("new request err: %v, param: %+v", err, slf)
		return
	}

	return call(request, slf.headers, slf.basic)
}

// Call 调用.
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
		err = fmt.Errorf("new request err: %v, param: %+v", err, slf)
		return
	}

	request.Header.Set("Content-Type", slf.contentType)

	return call(request, slf.headers, slf.basic)
}

// Call 调用.
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
		err = fmt.Errorf("new request err: %v, param: %+v", err, slf)
		return
	}
	request.Header.Set("Content-Type", slf.contentType)

	return call(request, slf.headers, slf.basic)
}

// Call 调用.
func (slf *RequestPostFormWithFile) Call() (respBytes []byte, err error) {
	var (
		buff   bytes.Buffer
		writer = multipart.NewWriter(&buff)
	)

	// 添加文件
	for k, v := range slf.Files {
		w, _ := writer.CreateFormFile(k, k)
		_, _ = w.Write(v)
	}
	// 添加表单参数
	for k, v := range slf.Body {
		_ = writer.WriteField(k, v)
	}
	_ = writer.Close()

	// 构建请求
	request, err := http.NewRequest(slf.method, slf.url, &buff)
	if err != nil {
		err = fmt.Errorf("new request err: %v, param: %+v", err, slf)
		return
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	return call(request, slf.headers, slf.basic)
}

// call 发送请求.
func call(request *http.Request, headers map[string]string, basic *BasicAuth) (respBytes []byte, err error) {
	var (
		resp *http.Response
	)

	// 防止Tcp复用
	request.Close = true
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	// 添加basic
	if basic != nil {
		request.SetBasicAuth(basic.User, basic.Pwd)
	}

	client := &http.Client{}
	client.Timeout = 30 * time.Second
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
		err = fmt.Errorf("readAll body err: %v, body: %s", err, resp.Body)
		return
	}

	return
}
