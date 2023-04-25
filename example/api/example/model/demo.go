package model

import (
	"fly/internal/domain/sqldb"
	"fly/internal/domain/types"
	"github.com/shopspring/decimal"
)

type (
	DemoCreateReq struct {
		Name   string          `json:"name" validate:"nonzero"`
		Amount decimal.Decimal `json:"amount" validate:"nonzero"`
		IsFree int             `json:"isFree" validate:"nonzero,min=1,max=2"`
		Remark string          `json:"remark" validate:"nonzero"`
	}

	DemoRecordReq struct {
		Name string `json:"name"`
		types.Pager
	}
	DemoRecordItem struct {
		Name   string          `json:"name"`
		Amount decimal.Decimal `json:"amount"`
		IsFree int             `json:"isFree"`
		Remark string          `json:"remark"`
	}
	DemoRecordResp struct {
		List []*DemoRecordItem `json:"list"`
		types.Pager
	}

	DemoInfoReq struct {
		DemoID int64 `json:"demoID" validate:"nonzero"`
	}
	DemoInfoResp struct {
		sqldb.Demo
	}

	DemoUpdateReq struct {
		DemoID int64           `json:"demoID" validate:"nonzero"`
		Name   string          `json:"name"`
		Amount decimal.Decimal `json:"amount"`
		IsFree int             `json:"isFree" validate:"min=0,max=2"`
		Remark string          `json:"remark"`
	}

	DemoDeleteReq struct {
		DemoID int64 `json:"demoID" validate:"nonzero"`
	}
)
