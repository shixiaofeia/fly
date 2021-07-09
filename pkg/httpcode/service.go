package httpcode

import (
	"bytes"
	"encoding/json"
	constants "fly/internal/const"
	"fly/pkg/logging"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 请求方式
const (
	RequestGet     = "GET"
	RequestGetPost = "POST"
)

// CallApi 调用接口
func CallApi(url, method string, data interface{}, response interface{}, header map[string]string) (respBytes []byte, err error) {
	// 加载参数
	reader := bytes.NewReader([]byte{})
	if data != nil {
		dataByte, _ := json.Marshal(data)
		reader = bytes.NewReader(dataByte)
	}
	// 构建请求
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		logging.Log.Error(fmt.Sprintf("NewRequest err: %v, url: %s", err, url))
		return
	}
	// 防止Tcp复用
	request.Close = true
	if method == RequestGetPost {
		request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}
	client := &http.Client{}
	client.Timeout = constants.CallInnerTimeOut
	resp, err := client.Do(request)
	if err != nil {
		logging.Log.Error(fmt.Sprintf("CallInnerWithHeader err: %v, url: %s", err, url))
		return
	}

	defer request.Body.Close()
	defer resp.Body.Close()

	// 读取返回
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Log.Error(fmt.Sprintf("CallInnerWithHeader ReadAll err: %v, url: %s, body: %v", err, url, resp.Body))
		return
	}
	// 需要结果则赋值结构体
	if response != nil {
		_ = json.Unmarshal(respBytes, response)
	}
	return
}
