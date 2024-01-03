package model

import (
	"fly/internal/domain/sqldb"
	"fly/internal/domain/types"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type (
	DemoCreateReq struct {
		Name   string            `json:"name" binding:"required"`
		Amount decimal.Decimal   `json:"amount" binding:"required"`
		IsFree int               `json:"isFree" binding:"required,min=1,max=2"`
		Remark string            `json:"remark" binding:"required"`
		Items  sqldb.StructItem  `json:"items" binding:"required"`
		Info   datatypes.JSONMap `json:"info" binding:"required"`
	}

	DemoRecordReq struct {
		Name string `json:"name"`
		types.Pager
	}
	DemoRecordItem struct {
		ID     int64           `json:"id"`
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
		DemoID int64 `json:"demoID" binding:"required"`
	}
	DemoInfoResp struct {
		sqldb.Demo
	}

	DemoUpdateReq struct {
		DemoID int64           `json:"demoID" binding:"required"`
		Name   string          `json:"name"`
		Amount decimal.Decimal `json:"amount"`
		IsFree int             `json:"isFree" binding:"min=0,max=2"`
		Remark string          `json:"remark"`
	}

	DemoDeleteReq struct {
		DemoID int64 `json:"demoID" binding:"required"`
	}
)
