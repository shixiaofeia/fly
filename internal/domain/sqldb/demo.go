package sqldb

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	StructItem struct {
		Key string
		Val int64
	}
	Demo struct {
		ID        int64             `gorm:"column:id;primaryKey;not null;autoIncrement" json:"id"`
		Name      string            `gorm:"column:name;type:varchar(20);not null;comment:名字" json:"name"`
		Amount    decimal.Decimal   `gorm:"column:amount;type:decimal(30,2);not null;comment:数量"  json:"amount"`
		IsFree    int               `gorm:"column:is_free;type:tinyint(1);not null;default:2;comment:是否免费 1是2否" json:"isFree"`
		Remark    string            `gorm:"column:remark;type:text;not null;comment:备注" json:"remark"`
		Items     StructItem        `gorm:"column:items;type:json;comment:扩展信息" json:"items"`
		Info      datatypes.JSONMap `gorm:"column:info;type:json;comment:扩展信息" json:"info"`
		CreatedAt int64             `gorm:"column:created_at;type:bigint(20);not null;comment:创建时间" json:"-"`
		UpdatedAt int64             `gorm:"column:updated_at;type:bigint(20);not null;comment:更新时间" json:"-"`
		DeletedAt gorm.DeletedAt    `gorm:"column:deleted_at;comment:删除时间;index" json:"-"`
	}
)

// Value 存储.
func (i StructItem) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	return string(b), err
}

// Scan 读取.
func (i *StructItem) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), i)
}
