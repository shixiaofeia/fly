package ckdb

import (
	"encoding/json"
	"fly/internal/domain/types"
	"fmt"
	"github.com/mailru/dbr"
	"github.com/shopspring/decimal"
)

type (
	DemoSearch struct {
		session *dbr.Session
		Demo
		types.Pager
		Order string
	}
)

// TableName 指定表名.
func (Demo) TableName() string {
	return "demo"
}

func (slf Demo) MarshalBinary() ([]byte, error) {
	return json.Marshal(slf)
}

// NewDemoSearch 实例Demo操作.
func NewDemoSearch(session *dbr.Session) *DemoSearch {
	return &DemoSearch{session: session}
}

// CreateTable 建表.
func (slf *DemoSearch) CreateTable() error {
	_, err := slf.session.Exec(fmt.Sprintf("CREATE TABLE If Not Exists %s  ("+
		"`id` String, "+
		"`user_id` String, "+
		"`order_id` String, "+
		"`currency` String, "+
		"`from_wallet` Int64, "+
		"`to_wallet` Int64, "+
		"`op_type` Int64, "+
		"`amount` Decimal128(18), "+
		"`c_time` Int64, "+
		"`is_ok` UInt8"+
		") ENGINE = MergeTree()  PARTITION BY toYYYYMM(toDateTime(create_time))  "+
		"ORDER BY (create_time)  SETTINGS index_granularity = 8192;", slf.TableName()))
	if err != nil {
		return fmt.Errorf("create table err: %v", err)
	}

	return nil
}

// BatchCreate 批量插入.
func (slf *DemoSearch) BatchCreate(records []*Demo) error {
	add := slf.session.InsertInto(slf.TableName()).Columns("id", "user_id", "order_id", "currency",
		"from_wallet", "to_wallet", "op_type", "amount", "c_time", "is_ok")
	for _, v := range records {
		add.Record(v)
	}
	_, err := add.Exec()

	return err
}

func (slf *DemoSearch) SetPager(page, size int) *DemoSearch {
	slf.PageNum, slf.PageSize = page, size
	return slf
}

func (slf *DemoSearch) SetOrder(order string) *DemoSearch {
	slf.Order = order

	return slf
}

func (slf *DemoSearch) SetID(id string) *DemoSearch {
	slf.ID = id
	return slf
}

func (slf *DemoSearch) SetUserId(userId string) *DemoSearch {
	slf.UserId = userId
	return slf
}

func (slf *DemoSearch) SetOrderId(orderId string) *DemoSearch {
	slf.OrderId = orderId
	return slf
}

func (slf *DemoSearch) SetCurrency(currency string) *DemoSearch {
	slf.Currency = currency
	return slf
}

func (slf *DemoSearch) SetFromWallet(fromWallet int) *DemoSearch {
	slf.FromWallet = fromWallet
	return slf
}

func (slf *DemoSearch) SetToWallet(toWallet int) *DemoSearch {
	slf.ToWallet = toWallet
	return slf
}

func (slf *DemoSearch) SetOpType(opType int) *DemoSearch {
	slf.OpType = opType
	return slf
}

func (slf *DemoSearch) SetAmount(amount decimal.Decimal) *DemoSearch {
	slf.Amount = amount
	return slf
}

func (slf *DemoSearch) SetCTime(cTime int64) *DemoSearch {
	slf.CTime = cTime
	return slf
}

func (slf *DemoSearch) SetIsOk(isOk bool) *DemoSearch {
	slf.IsOk = isOk
	return slf
}

// build 构建条件.
func (slf *DemoSearch) build(field string) *dbr.SelectBuilder {
	var builder = slf.session.Select(field).From(slf.TableName())

	if slf.ID != "" {
		builder = builder.Where("id = ?", slf.ID)
	}

	if slf.UserId != "" {
		builder = builder.Where("user_id = ?", slf.UserId)
	}

	if slf.OrderId != "" {
		builder = builder.Where("order_id = ?", slf.OrderId)
	}

	if slf.Currency != "" {
		builder = builder.Where("currency = ?", slf.Currency)
	}

	if slf.FromWallet > 0 {
		builder = builder.Where("from_wallet = ?", slf.FromWallet)
	}

	if slf.ToWallet > 0 {
		builder = builder.Where("to_wallet = ?", slf.ToWallet)
	}

	if slf.OpType > 0 {
		builder = builder.Where("op_type = ?", slf.OpType)
	}

	if slf.Amount.IsPositive() {
		builder = builder.Where("amount = ?", slf.Amount)
	}

	if slf.CTime > 0 {
		builder = builder.Where("c_time = ?", slf.CTime)
	}

	return builder
}

// Count 统计查询.
func (slf *DemoSearch) Count() (int64, error) {
	var (
		total int64
	)

	if _, err := slf.build("count(*)").Load(&total); err != nil {
		return total, fmt.Errorf("count err: %v", err)
	}

	return total, nil
}

// Find 多条查询.
func (slf *DemoSearch) Find() ([]*Demo, int64, error) {
	var (
		records = make([]*Demo, 0)
		total   int64
	)

	if _, err := slf.build("count(*)").Load(&total); err != nil {
		return records, total, fmt.Errorf("count err: %v", err)
	}

	builder := slf.build("*")
	if slf.PageNum > 0 && slf.PageSize > 0 {
		builder = builder.Offset(uint64((slf.PageNum - 1) * slf.PageSize)).
			Limit(uint64(slf.PageSize))
	}
	if slf.Order != "" {
		builder = builder.OrderBy(slf.Order)
	}
	if _, err := builder.Load(&records); err != nil {
		return records, total, fmt.Errorf("find err: %v", err)
	}

	return records, total, nil
}

// FindNotTotal 多条查询&不查询总数.
func (slf *DemoSearch) FindNotTotal() ([]*Demo, error) {
	var (
		records = make([]*Demo, 0)
	)

	builder := slf.build("*")
	if slf.PageNum > 0 && slf.PageSize > 0 {
		builder = builder.Offset(uint64((slf.PageNum - 1) * slf.PageSize)).
			Limit(uint64(slf.PageSize))
	}
	if slf.Order != "" {
		builder = builder.OrderBy(slf.Order)
	}
	if _, err := builder.Load(&records); err != nil {
		return records, fmt.Errorf("find err: %v", err)
	}

	return records, nil
}

// First 单条查询.
func (slf *DemoSearch) First() (*Demo, error) {
	var (
		recordM = new(Demo)
	)

	if _, err := slf.build("*").Load(recordM); err != nil {
		return recordM, fmt.Errorf("first err: %v", err)
	}

	return recordM, nil
}

// DeleteByID 删除.
func (slf *DemoSearch) DeleteByID(id string) error {
	if _, err := slf.session.DeleteFrom(slf.TableName()).Where("id = ?", id).Exec(); err != nil {
		return fmt.Errorf("delete err: %v", err)
	}

	return nil
}
