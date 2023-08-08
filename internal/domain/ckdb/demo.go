package ckdb

import "github.com/shopspring/decimal"

type (
	Demo struct {
		ID         string
		UserId     string
		OrderId    string
		Currency   string
		FromWallet int
		ToWallet   int
		OpType     int
		Amount     decimal.Decimal
		CTime      int64
		IsOk       bool
	}
)
