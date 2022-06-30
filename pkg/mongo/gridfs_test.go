package mongo

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestNewGridFS(t *testing.T) {
	var (
		err  error
		conf = Conf{
			Host:        "127.0.0.1",
			Port:        "27017",
			DataBase:    "file",
			MaxPoolSize: 10,
		}
		buff        = new(bytes.Buffer)
		contentType = "image/png"
		fileName    = "test.png"
	)

	if err = Init(conf); err != nil {
		t.Fatalf("init err: %v", err)
		return
	}

	cli := NewGridFS("fs")

	// 准备文件
	fd, _ := os.Open(fmt.Sprintf("./%s", fileName))
	defer fd.Close()
	_, _ = io.Copy(buff, fd)

	// 创建
	fid, err := cli.Create(buff.Bytes(), fileName, contentType)
	if err != nil {
		t.Fatalf("create err: %v", err)
		return
	}
	t.Logf("create file id: %v", fid)

	// 查询
	fileContent, err := cli.OpenById(fid)
	if err != nil {
		t.Fatalf("open err: %v", err)
		return
	}
	t.Logf("image base64: %s", base64.StdEncoding.EncodeToString(fileContent))

	// 移除
	if err = cli.Remove(fid); err != nil {
		t.Fatalf("remove err: %v", err)
		return
	}

	t.Log("test success")
}
