package sqldb

import (
	"fly/internal/constvar"
	"fly/pkg/mysql"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
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

// NewDemoSearch 实例Demo操作.
func NewDemoSearch(db *gorm.DB) *DemoSearch {
	if db == nil {
		db = mysql.NewWriteDB()
	}

	return &DemoSearch{Orm: db}
}

func (slf *DemoSearch) SetOrm(tx *gorm.DB) *DemoSearch {
	slf.Orm = tx

	return slf
}

func (slf *DemoSearch) SetPager(page, size int) *DemoSearch {
	slf.Page, slf.Size = page, size

	return slf
}

func (slf *DemoSearch) SetOrder(order string) *DemoSearch {
	slf.Order = order

	return slf
}

func (slf *DemoSearch) SetId(id uint) *DemoSearch {
	slf.Id = id
	return slf
}

func (slf *DemoSearch) SetName(name string) *DemoSearch {
	slf.Name = name
	return slf
}

func (slf *DemoSearch) SetAmount(amount decimal.Decimal) *DemoSearch {
	slf.Amount = amount
	return slf
}

func (slf *DemoSearch) SetIsFree(isFree int) *DemoSearch {
	slf.IsFree = isFree
	return slf
}

func (slf *DemoSearch) SetRemark(remark string) *DemoSearch {
	slf.Remark = remark
	return slf
}

func (slf *DemoSearch) SetCreateTime(createTime int64) *DemoSearch {
	slf.CreateTime = createTime
	return slf
}

func (slf *DemoSearch) SetUpdateTime(updateTime int64) *DemoSearch {
	slf.UpdateTime = updateTime
	return slf
}

func (slf *DemoSearch) SetDeleteTime(deleteTime int64) *DemoSearch {
	slf.DeleteTime = deleteTime
	return slf
}

// CreateTable 创建表.
func (slf *DemoSearch) CreateTable() {
	_ = slf.Orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Demo'").AutoMigrate(&Demo{})
}

// Create 单条插入
func (slf *DemoSearch) Create(recordM *Demo) (uint, error) {
	if err := slf.Orm.Create(recordM).Error; err != nil {
		return 0, fmt.Errorf("create recordM err: %v, recordM: %+v", err, recordM)
	}

	return recordM.Id, nil
}

// CreateBatch 批量插入
func (slf *DemoSearch) CreateBatch(records []*Demo) error {
	if err := slf.Orm.Create(&records).Error; err != nil {
		return fmt.Errorf("create records err: %v, records: %+v", err, records)
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

	if slf.Amount.IsPositive() {
		db = db.Where("amount = ?", slf.Amount)
	}

	if slf.IsFree > 0 {
		db = db.Where("is_free = ?", slf.IsFree)
	}

	if slf.Remark != "" {
		db = db.Where("remark = ?", slf.Remark)
	}

	if slf.CreateTime > 0 {
		db = db.Where("create_time = ?", slf.CreateTime)
	}

	if slf.UpdateTime > 0 {
		db = db.Where("update_time = ?", slf.UpdateTime)
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

// FindByQuery 指定条件查询.
func (slf *DemoSearch) FindByQuery(query string, args []interface{}) ([]*Demo, int64, error) {
	var (
		records = make([]*Demo, 0)
		total   int64
		db      = slf.Orm.Table(slf.TableName()).Where(query, args...)
	)

	if err := db.Count(&total).Error; err != nil {
		return records, total, fmt.Errorf("find by query count err: %v", err)
	}
	if slf.Page*slf.Size > 0 {
		db = db.Offset((slf.Page - 1) * slf.Size).Limit(slf.Size)
	}
	if err := db.Find(&records).Error; err != nil {
		return records, total, fmt.Errorf("find by query records err: %v, query: %s, args: %+v", err, query, args)
	}

	return records, total, nil
}

// FindByQueryNotTotal 指定条件查询&不查询总条数.
func (slf *DemoSearch) FindByQueryNotTotal(query string, args []interface{}) ([]*Demo, error) {
	var (
		records = make([]*Demo, 0)
		db      = slf.Orm.Table(slf.TableName()).Where(query, args...)
	)

	if slf.Page*slf.Size > 0 {
		db = db.Offset((slf.Page - 1) * slf.Size).Limit(slf.Size)
	}
	if err := db.Find(&records).Error; err != nil {
		return records, fmt.Errorf("find by query records err: %v, query: %s, args: %+v", err, query, args)
	}

	return records, nil
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
		db = slf.build()
	)

	if err := db.Updates(upMap).Error; err != nil {
		return fmt.Errorf("update err: %v", err)
	}

	return nil
}

// UpdateByQuery 指定条件更新.
func (slf *DemoSearch) UpdateByQuery(query string, args []interface{}, upMap map[string]interface{}) error {
	var (
		db = slf.Orm.Table(slf.TableName()).Where(query, args...)
	)

	if err := db.Updates(upMap).Error; err != nil {
		return fmt.Errorf("update by query err: %v, query: %s, args: %+v, upMap: %+v", err, query, args, upMap)
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
