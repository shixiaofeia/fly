package sqldb

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Domo struct {
	Id         uint32          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name       string          `gorm:"column:name;type:varchar(20);default:'';comment:名字" json:"name"`
	Amount     decimal.Decimal `gorm:"column:amount;type:decimal(30,2);default:0;comment:数量"  json:"amount"`
	IsFree     int             `gorm:"column:is_free;type:tinyint(1);default:2;comment:是否免费 1是2否" json:"isFree"`
	Remark     string          `gorm:"column:remark;type:text;comment:备注" json:"remark"`
	CreateTime int64           `gorm:"column:create_time;type:bigint(20);default:0;comment:创建时间" json:"createTime"`
	UpdateTime int64           `gorm:"column:update_time;type:bigint(20);default:0;comment:更新时间" json:"updateTime"`
	DeleteTime int64           `gorm:"column:delete_time;type:bigint(20);default:0;comment:删除时间" json:"-"`
}

// TableName 指定表名.
func (d *Domo) TableName() string {
	return "demo"
}

// CreateTable 创建表.
func (d *Domo) CreateTable(db *gorm.DB) {
	_ = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='demo'").AutoMigrate(&Domo{})
}
