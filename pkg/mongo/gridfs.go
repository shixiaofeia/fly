package mongo

import (
	"bytes"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io"
	"time"
)

// OpenById 通过id打开文件.
func (slf *GridFS) OpenById(id interface{}) ([]byte, error) {
	var (
		buff = new(bytes.Buffer)
	)
	switch id.(type) {
	case string:
		id = bson.ObjectIdHex(id.(string))
	}
	file, err := slf.client.OpenId(id)
	if err != nil {
		return nil, fmt.Errorf("open gridfs err: %v", err)
	}
	_, _ = io.Copy(buff, file)
	_ = file.Close()
	return buff.Bytes(), nil
}

// Create 创建文件.
func (slf *GridFS) Create(fileContent []byte, fileName, contentType string) (string, error) {
	writer, err := slf.client.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("create gridfs writer err: %v", err)
	}
	defer writer.Close()

	if _, err = writer.Write(fileContent); err != nil {
		return "", fmt.Errorf("witre contennt err: %v", err)
	}
	writer.SetContentType(contentType)
	writer.SetUploadDate(time.Now().Add(8 * time.Hour)) // 解决时区差的问题
	id, ok := writer.Id().(bson.ObjectId)
	if ok {
		return id.Hex(), nil
	}
	return "", fmt.Errorf("id type err, val: %v", id)
}

// Remove 移除.
func (slf *GridFS) Remove(id interface{}) error {
	switch id.(type) {
	case string:
		id = bson.ObjectIdHex(id.(string))
	}
	if err := slf.client.RemoveId(id); err != nil {
		return fmt.Errorf("gridfs remove err: %v", err)
	}
	return nil
}
