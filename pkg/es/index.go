package es

import (
	"context"
	"github.com/olivere/elastic/v6"
)

var client *elastic.Client

// Init 初始化es
func Init(address ...string) (err error) {
	if client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(address...)); err != nil {
		return
	}
	if _, _, err = client.Ping(address[0]).Do(context.Background()); err != nil {
		return
	}
	return
}

// NewClient
func NewClient() *elastic.Client {
	return client
}
