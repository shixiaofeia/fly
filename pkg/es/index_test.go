package es

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/olivere/elastic/v6"
	"github.com/shopspring/decimal"
)

type Fly struct {
	UserId     string          `json:"userId"`
	Remark     string          `json:"remark"`
	Amount     decimal.Decimal `json:"amount"`
	CreateTime int64           `json:"createTime"`
}

func TestNewClient(t *testing.T) {
	var (
		indexName = "fly"
		indexType = "_doc"
	)
	if err := Init("http://127.0.0.1:9200"); err != nil {
		t.Error(err.Error())
	}
	// index exist
	b, err := NewClient().IndexExists(indexName).Do(context.Background())
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(b)
	if !b {
		// create index
		if _, err = NewClient().CreateIndex(indexName).BodyString(getMapping()).Do(context.Background()); err != nil {
			t.Error(err)
		}
	}
	// add data
	flyM := Fly{UserId: "10010", Remark: "备注", Amount: decimal.NewFromInt(1), CreateTime: time.Now().Unix()}
	if _, err = NewClient().Index().Index(indexName).Type(indexType).BodyJson(flyM).Do(context.Background()); err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Second)
	// query
	query := elastic.NewTermQuery("userId", "10010")
	t.Log(query.Source())
	searchRes, err := NewClient().Search(indexName).Query(query).Do(context.Background())
	if err != nil {
		t.Error(err)
	}
	if len(searchRes.Hits.Hits) > 0 {
		for _, v := range searchRes.Hits.Hits {
			var f Fly
			_ = json.Unmarshal(*v.Source, &f)
			fmt.Printf("data: %+v", f)
		}
	}
	// del index
	if _, err = NewClient().DeleteIndex(indexName).Do(context.Background()); err != nil {
		t.Error(err)
	}
}

// getMapping 获取mapping.
func getMapping() string {
	return `
{
 "mappings": {
   "_doc": {
     "properties": {
       "userId": {
         "type": "keyword"
       },
 		"remark": {
         "type": "text"
       },
       "amount": {
         "type": "scaled_float",
         "scaling_factor": 100
       },
		"createTime": {
         "type": "long"
       }
     }
   }
 }
}
`
}
