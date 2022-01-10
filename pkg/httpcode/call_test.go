package httpcode

import (
	"testing"
)

func TestCallApi(t *testing.T) {
	var (
		rg   = NewCallByGet("http://baidu.com")
		rj   = NewCall(RequestPostJson, "", map[string]interface{}{}, nil)
		rf   = NewCall(RequestPostForm, "", map[string]interface{}{}, nil)
		resp []byte
		err  error
	)
	resp, err = rg.Call()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(string(resp))

	resp, err = rj.Call()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(string(resp))

	resp, err = rf.Call()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(string(resp))
}
