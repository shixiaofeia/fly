package main

import (
	"fly/internal/domain/sqldb"
	"fmt"
	"reflect"
	"strings"
)

func main() {
	var (
		obj       = sqldb.Demo{}
		objName   string
		buildFunc string
		setFunc   string
	)

	filedT := reflect.TypeOf(obj)
	objName = filedT.Name()
	filedNum := filedT.NumField()

	for i := 0; i < filedNum; i++ {
		val := filedT.Field(i)
		fName := val.Name
		gName := ColumnName(val.Tag.Get("gorm"))
		fType := val.Type.Name()
		//fmt.Println(fName, fType, gName)
		if v := getWhere(fName, fType, gName); v != "" {
			buildFunc += v
		}
		if v := getSetFunc(objName, fName, fType); v != "" {
			setFunc += v
		}
	}

	buildFunc = fmt.Sprintf(`
// build 构建条件.
func (slf *%sSearch) build() *gorm.DB {
	var db = slf.Orm.Table(slf.TableName()).Model(%s{})
	%s
	if slf.Order != "" {
		db = db.Order(slf.Order)
	}

	return db
}
`, objName, objName, buildFunc)

	template := strings.Replace(getFormat(objName), "buildFunc", buildFunc, 1)
	template = strings.Replace(template, "setFunc", setFunc, 1)

	fmt.Println(template)
}

func ColumnName(s string) string {
	for _, name := range strings.Split(s, ";") {
		if strings.Index(name, "column") == -1 {
			continue
		}
		return strings.Replace(name, "column:", "", 1)
	}

	return ""
}

func getWhere(fName, fType, gName string) string {
	if strings.Contains(fType, "int") {
		return fmt.Sprintf(`
	if slf.%s > 0 {
		db = db.Where("%s = ?", slf.%s)
	}			
		`, fName, gName, fName)
	} else if strings.Contains(fType, "string") {
		return fmt.Sprintf(`
	if slf.%s != "" {
		db = db.Where("%s = ?", slf.%s)
	}			
		`, fName, gName, fName)
	} else if strings.Contains(fType, "Decimal") {
		return fmt.Sprintf(`
	if slf.%s.IsPositive() {
		db = db.Where("%s = ?", slf.%s)
	}			
		`, fName, gName, fName)
	}

	return ""
}

func getSetFunc(objName, fName, fType string) string {
	if fName == "DeletedAt" {
		return ""
	}
	var param = strings.ToLower(fName[:1]) + fName[1:]
	if fName == "ID" {
		param = strings.ToLower(fName)
	}
	if fType == "Decimal" {
		fType = "decimal.Decimal"
	}
	str := fmt.Sprintf(`
func (slf *%sSearch) Set%s(%s %s) *%sSearch {
	slf.%s = %s
	return slf
}
`, objName, fName, param, fType, objName, fName, param)

	return str
}

func getFormat(objName string) string {
	tableName := strings.ToLower(objName)
	str := `

type (
	DemoSearch struct {
		Orm *gorm.DB
		Demo
		types.Pager
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
	slf.PageNum, slf.PageSize = page, size
	return slf
}

func (slf *DemoSearch) SetOrder(order string) *DemoSearch {
	slf.Order = order

	return slf
}
setFunc

// CreateTable 创建表.
func (slf *DemoSearch) CreateTable() {
	_ = slf.Orm.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Demo'").AutoMigrate(&Demo{})
}

// Create 单条插入
func (slf *DemoSearch) Create(recordM *Demo) (int64, error) {
	if err := slf.Orm.Create(recordM).Error; err != nil {
		return 0, fmt.Errorf("create recordM err: %v, recordM: %+v", err, recordM)
	}

	return recordM.ID, nil
}

// CreateBatch 批量插入
func (slf *DemoSearch) CreateBatch(records []*Demo) error {
	if err := slf.Orm.Create(&records).Error; err != nil {
		return fmt.Errorf("create records err: %v, records: %+v", err, records)
	}

	return nil
}

buildFunc

// Count 统计查询.
func (slf *DemoSearch) Count() (int64, error) {
	var (
		total   int64
		db      = slf.build()
	)

	if err := db.Count(&total).Error; err != nil {
		return total, fmt.Errorf("count err: %v", err)
	}

	return total, nil
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
	if slf.PageNum*slf.PageSize > 0 {
		db = db.Offset((slf.PageNum - 1) * slf.PageSize).Limit(slf.PageSize)
	}
	if err := db.Find(&records).Error; err != nil {
		return records, total, fmt.Errorf("find records err: %v", err)
	}

	return records, total, nil
}

// FindNotTotal 多条查询&不查询总数.
func (slf *DemoSearch) FindNotTotal() ([]*Demo, error) {
	var (
		records = make([]*Demo, 0)
		db      = slf.build()
	)

	if slf.PageNum*slf.PageSize > 0 {
		db = db.Offset((slf.PageNum - 1) * slf.PageSize).Limit(slf.PageSize)
	}
	if err := db.Find(&records).Error; err != nil {
		return records, fmt.Errorf("find records err: %v", err)
	}

	return records, nil
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
	if slf.PageNum*slf.PageSize > 0 {
		db = db.Offset((slf.PageNum - 1) * slf.PageSize).Limit(slf.PageSize)
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

	if slf.PageNum*slf.PageSize > 0 {
		db = db.Offset((slf.PageNum - 1) * slf.PageSize).Limit(slf.PageSize)
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
		return recordM, err
	}

	return recordM, nil
}

// UpdateByMap 更新.
func (slf *DemoSearch) UpdateByMap(upMap map[string]interface{}) error {
	var (
		db = slf.build()
	)

	if err := db.Updates(upMap).Error; err != nil {
		return fmt.Errorf("update by map err: %v, upMap: %+v", err, upMap)
	}

	return nil
}

// UpdateByStruct 更新.
func (slf *DemoSearch) UpdateByStruct(recordM *Demo) error {
	var (
		db = slf.build()
	)

	if err := db.Updates(recordM).Error; err != nil {
		return fmt.Errorf("update by struct err: %v, record: %+v", err, recordM)
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

// SaveByStruct 更新整个结构体(带ID为更新否则新增).
func (slf *DemoSearch) SaveByStruct(recordM *Demo) error {
	var (
		db = slf.build()
	)

	if err := db.Save(recordM).Error; err != nil {
		return fmt.Errorf("save by struct err: %v, record: %+v", err, recordM)
	}

	return nil
}

// Delete 删除.
func (slf *DemoSearch) Delete() error {
	var (
		db = slf.build()
	)

	if err := db.Delete(&Demo{}).Error; err != nil {
		return fmt.Errorf("delete err: %v", err)
	}

	return nil
}

`
	str = strings.ReplaceAll(str, "Demo", objName)
	str = strings.ReplaceAll(str, "demo", tableName)
	return str
}
