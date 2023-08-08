package main

import (
	"fly/internal/domain/ckdb"
	"fmt"
	"reflect"
	"strings"
)

// var filterMap = map[string]string{
// 	"ID":  "id",
// 	"IP":  "ip",
// 	"URL": "url",
// }

func main() {
	var (
		obj             = ckdb.Demo{}
		objName         string
		buildFunc       string
		setFunc         string
		createTableFunc string
		fieldList       []string
	)

	filedT := reflect.TypeOf(obj)
	objName = filedT.Name()
	filedNum := filedT.NumField()
	tableName := strings.ToLower(objName)

	for i := 0; i < filedNum; i++ {
		val := filedT.Field(i)
		fName := val.Name
		gName := ToSnakeCase(val.Name)
		fType := val.Type.Name()
		// fmt.Println(fName, fType, gName)
		if v := getWhere(fName, fType, gName); v != "" {
			buildFunc += v
		}
		if v := getSetFunc(objName, fName, fType); v != "" {
			setFunc += v
		}
		if v := getCreateTable(fType, gName); v != "" {
			createTableFunc += v
		}
		fieldList = append(fieldList, fmt.Sprintf(`"%s"`, gName))
	}

	buildFunc = fmt.Sprintf(`
// build 构建条件.
func (slf *%sSearch) build(field string) *dbr.SelectBuilder {
    var builder = slf.session.Select(field).From(slf.TableName())
	%s

	return builder
}
`, objName, buildFunc)
	createTableFunc = replaceObjName(getCreateTableFunc(createTableFunc), objName)

	createFunc := replaceObjName(getCreateFunc(fieldList), objName)
	template := strings.Replace(getFormat(objName, tableName), "buildFunc", buildFunc, 1)
	template = strings.Replace(template, "setFunc", setFunc, 1)
	template = strings.Replace(template, "createFunc", createFunc, 1)
	template = strings.Replace(template, "createTableFunc", createTableFunc, 1)

	fmt.Println(template)
}

func ToSnakeCase(s string) string {
	if s == "ID" {
		return "id"
	}
	var result strings.Builder
	for i, r := range s {
		if i > 0 && ('A' <= r && r <= 'Z') {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func getCreateTable(fType, gName string) string {
	if strings.Contains(fType, "int") {
		return fmt.Sprintf("`%s` Int64, ", gName)
	} else if strings.Contains(fType, "string") {
		return fmt.Sprintf("`%s` String, ", gName)
	} else if strings.Contains(fType, "Decimal") {
		return fmt.Sprintf("`%s` Decimal128(18), ", gName)
	} else if strings.Contains(fType, "bool") {
		return fmt.Sprintf("`%s` UInt8", gName)
	}

	return ""
}

func getCreateTableFunc(createTableFunc string) string {
	return fmt.Sprintf(`
// CreateTable 建表 todo 修改字段类型, 分区及排序字段.
func (slf *DemoSearch) CreateTable() error {
	_, err := slf.session.Exec(fmt.Sprintf("CREATE TABLE If Not Exists %%s  ("+
		"%s"+
		") ENGINE = MergeTree()  PARTITION BY toYYYYMM(toDateTime(create_time))  "+
		"ORDER BY (create_time)  SETTINGS index_granularity = 8192;", slf.TableName()))
	if err != nil {
		return fmt.Errorf("create table err: %%v", err)
	}

	return nil
}
`, createTableFunc)
}

func getCreateFunc(fieldList []string) string {
	str := fmt.Sprintf(`
// BatchCreate 批量插入.
func (slf *DemoSearch) BatchCreate(records []*Demo) error {
	add := slf.session.InsertInto(slf.TableName()).Columns(%s)
	for _, v := range records {
		add.Record(v)
	}
	_, err := add.Exec()

	return err
}
`, strings.Join(fieldList, ","))

	return str
}

func getWhere(fName, fType, gName string) string {
	if strings.Contains(fType, "int") {
		return fmt.Sprintf(`
	if slf.%s > 0 {
		builder = builder.Where("%s = ?", slf.%s)
	}			
		`, fName, gName, fName)
	} else if strings.Contains(fType, "string") {
		return fmt.Sprintf(`
	if slf.%s != "" {
		builder = builder.Where("%s = ?", slf.%s)
	}			
		`, fName, gName, fName)
	} else if strings.Contains(fType, "Decimal") {
		return fmt.Sprintf(`
	if slf.%s.IsPositive() {
		builder = builder.Where("%s = ?", slf.%s)
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

func getFormat(objName, tableName string) string {
	str := `

type (
	DemoSearch struct {
		session *dbr.Session
		Demo
		types.Pager
		Order string
	}
)

// TableName 指定表名 todo 修改表名.
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

createTableFunc

createFunc

func (slf *DemoSearch) SetPager(page, size int) *DemoSearch {
	slf.PageNum, slf.PageSize = page, size
	return slf
}

func (slf *DemoSearch) SetOrder(order string) *DemoSearch {
	slf.Order = order

	return slf
}
setFunc
buildFunc
// Count 统计查询.
func (slf *DemoSearch) Count() (int64, error) {
	var (
		total   int64
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

`
	str = strings.ReplaceAll(str, "Demo", objName)
	str = strings.ReplaceAll(str, "demo", tableName)
	return str
}

func replaceObjName(str, objName string) string {
	str = strings.ReplaceAll(str, "Demo", objName)
	return str
}
