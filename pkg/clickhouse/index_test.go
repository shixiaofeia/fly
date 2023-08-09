package clickhouse

import (
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/mailru/go-clickhouse"
	"github.com/shopspring/decimal"
)

type (
	asset struct {
		UserId     string
		OrderId    string
		Currency   string
		FromWallet int
		ToWallet   int
		OpType     int
		Amount     decimal.Decimal
		CTime      int64
	}
)

func TestInit(t *testing.T) {
	var err error
	// create database if not exists house;
	c := Config{Host: "127.0.0.1", Port: "8123", Database: "fly", User: "root", Pwd: "root"}
	if err = Init(c); err != nil {
		log.Fatal("init clickhouse err: " + err.Error())
	}

	sess := NewSession()
	_, err = sess.Exec("CREATE TABLE If Not Exists asset (`user_id` String, `order_id` String, `currency` String, " +
		"`from_wallet` UInt16, `to_wallet` UInt16, `op_type` UInt16, `amount` Decimal128(18), `c_time` UInt32) " +
		"ENGINE = MergeTree() PARTITION BY toYYYYMM(toDateTime(c_time)) ORDER BY (c_time) " +
		"SETTINGS index_granularity = 8192;")
	if err != nil {
		log.Fatal("create table err: " + err.Error())
	}

	add := sess.InsertInto("asset").Columns("user_id", "order_id", "currency", "from_wallet", "to_wallet",
		"op_type", "amount", "c_time")
	amount, _ := decimal.NewFromString("0.123456789123456789")
	for i := 0; i < 10000; i++ {
		add.Record(asset{
			UserId:     fmt.Sprintf("%d", i),
			OrderId:    fmt.Sprintf("%d", time.Now().UnixNano()),
			Currency:   "USD",
			FromWallet: 1,
			ToWallet:   1,
			OpType:     1,
			Amount:     amount.Add(decimal.NewFromInt(int64(i))),
			CTime:      time.Now().Unix(),
		})
	}
	res, err := add.Exec()
	log.Printf("res: %+v, err: %v", res, err)

	countQ := sess.SelectBySql("select count(*) as total from asset")
	var total int
	_, err = countQ.Load(&total)
	log.Printf("err: %v, total: %d", err, total)

	var items []struct {
		Amount     decimal.Decimal `json:"amount"`
		FromWallet int             `json:"fromWallet"`
	}
	query := sess.SelectBySql("select sum(amount) as amount, from_wallet from asset group by from_wallet")
	if _, err = query.Load(&items); err != nil {
		log.Fatal(err)
	}
	for _, item := range items {
		log.Printf("amount: %v, fromWallet: %v", item.Amount, item.FromWallet)
	}
}
