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
		setFunc += getSetFunc(objName, fName, fType)
	}

	buildFunc = fmt.Sprintf(`
// build 构建条件.
func (slf *%sSearch) build() *gorm.DB {
	var db = slf.Orm.Table(slf.TableName()).Where("delete_time=0")
	%s
	if slf.Order != "" {
		db = db.Order(slf.Order)
	}

	return db
}
`, objName, buildFunc)

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
	param := strings.ToLower(fName[:1]) + fName[1:]
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

setFunc

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

buildFunc

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
func (slf *DemoSearch) FindByQuery(query string, args []interface{}) ([]*Demo, error) {
	var (
		records = make([]*Demo, 0)
	)

	if err := slf.Orm.Table(slf.TableName()).Where(query, args...).Find(&records).Error; err != nil {
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
`
	str = strings.ReplaceAll(str, "Demo", objName)
	str = strings.ReplaceAll(str, "demo", tableName)
	return str
}
