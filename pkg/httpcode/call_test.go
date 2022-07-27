package httpcode

import (
	"testing"
)

func TestNewRequestGet(t *testing.T) {
	var (
		url     = ""
		factory = NewRequestGet(url)
	)
	factory.AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	resp, err := factory.Call()
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("resp: %s", resp)

}

func TestNewRequestPostJson(t *testing.T) {
	var (
		url  = ""
		data = map[string]interface{}{
			"name": "fly",
		}
		factory = NewRequestPostJson(url, data)
	)
	factory.AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	resp, err := factory.Call()
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("resp: %s", resp)
}

func TestNewRequestPostForm(t *testing.T) {
	var (
		url  = ""
		data = map[string]string{
			"name": "fly",
		}
		factory = NewRequestPostForm(url, data)
	)
	factory.AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	resp, err := factory.Call()
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("resp: %s", resp)
}

func TestNewRequestPostFormWithFile(t *testing.T) {
	var (
		url  = ""
		data = map[string]string{
			"name": "fly",
		}
		files = map[string][]byte{
			"content": []byte("this is content"),
		}
		factory = NewRequestPostFormWithFile(url, data, files)
	)
	factory.AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	resp, err := factory.Call()
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("resp: %s", resp)

}
