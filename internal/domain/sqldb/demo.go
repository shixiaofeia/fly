package sqldb

import (
	"fly/internal/constvar"
	"fly/pkg/mysql"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type (
	Demo struct {
		Id         uint            `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
		Name       string          `gorm:"column:name;type:varchar(20);default:'';comment:名字" json:"name"`
		Amount     decimal.Decimal `gorm:"column:amount;type:decimal(30,2);default:0;comment:数量"  json:"amount"`
		IsFree     int             `gorm:"column:is_free;type:tinyint(1);default:2;comment:是否免费 1是2否" json:"isFree"`
		Remark     string          `gorm:"column:remark;type:text;comment:备注" json:"remark"`
		CreateTime int64           `gorm:"column:create_time;type:bigint(20);default:0;comment:创建时间" json:"createTime"`
		UpdateTime int64           `gorm:"column:update_time;type:bigint(20);default:0;comment:更新时间" json:"updateTime"`
		DeleteTime int64           `gorm:"column:delete_time;type:bigint(20);default:0;comment:删除时间" json:"-"`
	}

	DemoSearch struct {
		Orm *gorm.DB
		Demo
		constvar.Pager
		Order string
	}
)

// TableName 指定表名.
func (Demo) TableName() string {
	return "demo"
}

// NewDemoSearch 实例demo操作.
func NewDemoSearch(db *gorm.DB) *DemoSearch {
	if db == nil {
		db = mysql.NewWriteDB()
	}
	return &DemoSearch{Orm: db}
}

// CreateTable 创建表.
func (slf *DemoSearch) CreateTable() {
	_ = slf.Orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='demo'").AutoMigrate(&Demo{})
}

// Create 单条插入
func (slf *DemoSearch) Create(recordM *Demo) (uint, error) {
	if err := slf.Orm.Create(recordM).Error; err != nil {
		return 0, fmt.Errorf("create recordM err: %v", err)
	}

	return recordM.Id, nil
}

// CreateBatch 批量插入
func (slf *DemoSearch) CreateBatch(records []*Demo) error {
	if err := slf.Orm.Create(&records).Error; err != nil {
		return fmt.Errorf("create records err: %v", err)
	}

	return nil
}

// build 构建条件.
func (slf *DemoSearch) build() *gorm.DB {
	var db = slf.Orm.Table(slf.TableName()).Where("delete_time=0")

	if slf.Id > 0 {
		db = db.Where("id = ?", slf.Id)
	}
	if slf.Name != "" {
		db = db.Where("name like ?", "%"+slf.Name+"%")
	}
	if slf.IsFree > 0 {
		db = db.Where("is_free=?", slf.IsFree)
	}
	if slf.Order != "" {
		db = db.Order(slf.Order)
	}

	return db
}

// Find 多条查询.
func (slf *DemoSearch) Find() ([]*Demo, int64, error) {
	var (
		records = make([]*Demo, 0)
		total   int64
		db      = slf.build()
	)

	if err := db.Count(&total).Error; err != nil {
		return records, total, fmt.Errorf("find count err: %v", err)
	}
	if slf.Page*slf.Size > 0 {
		db = db.Offset((slf.Page - 1) * slf.Size).Limit(slf.Size)
	}
	if err := db.Find(&records).Error; err != nil {
		return records, total, fmt.Errorf("find records err: %v", err)
	}

	return records, total, nil
}

// First 单条查询.
func (slf *DemoSearch) First() (*Demo, error) {
	var (
		recordM = new(Demo)
		db      = slf.build()
	)

	if err := db.First(recordM).Error; err != nil {
		return recordM, fmt.Errorf("frist recordM err: %v", err)
	}

	return recordM, nil
}

// Update 更新.
func (slf *DemoSearch) Update(upMap map[string]interface{}) error {
	var (
		db = slf.build().Debug()
	)

	if err := db.Updates(upMap).Error; err != nil {
		return fmt.Errorf("update err: %v", err)
	}

	return nil
}

// Delete 删除.
func (slf *DemoSearch) Delete() error {
	var (
		db = slf.build()
	)

	if err := db.Updates(map[string]interface{}{"delete_time": time.Now().Unix()}).Error; err != nil {
		return fmt.Errorf("delete err: %v", err)
	}

	return nil
}
