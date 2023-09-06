package es

import (
	"context"
	"time"

	"github.com/olivere/elastic/v6"
)

var client *elastic.Client

// Init 初始化es.
func Init(address ...string) (err error) {
	if client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(address...)); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, _, err = client.Ping(address[0]).Do(ctx); err != nil {
		return
	}

	return
}

// NewClient 获取client
func NewClient() *elastic.Client {
	return client
}
