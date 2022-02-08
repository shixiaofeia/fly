package httpcode

import (
	"testing"
)

func TestCallApi(t *testing.T) {

	var (
		url  = "http://baidu.com"
		rg   = NewRequestGet(url)
		rj   = NewRequestPostJson(url, map[string]interface{}{})
		rf   = NewRequestPostForm(url, map[string]string{})
		resp []byte
		err  error
	)
	rg.AddHeader("t", "1")
	resp, err = rg.Call()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(string(resp))

	rj.AddHeader("t", "2")
	resp, err = rj.Call()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(string(resp))

	rf.AddHeader("t", "3")
	resp, err = rf.Call()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(string(resp))
}
