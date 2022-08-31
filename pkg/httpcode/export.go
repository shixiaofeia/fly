package httpcode

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExportExcel 导出excel.
func ExportExcel(thead []string, tbody interface{}) (r *bytes.Buffer, err error) {
	// make sure 'tbody' is a Slice
	valTbody := reflect.ValueOf(tbody)
	if valTbody.Kind() != reflect.Slice {
		err = fmt.Errorf("tbody not slice")
		return
	}
	var (
		// it's a slice, so open up its values
		n     = valTbody.Len()
		sheet = "Sheet1"
	)

	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(sheet)

	// write thead
	for k, v := range thead {
		var axis string
		axis, err = excelize.CoordinatesToCellName(k+1, 1)
		if err != nil {
			err = fmt.Errorf("thead to cell name error: %s", err)
			return
		}
		_ = xlsx.SetCellValue(sheet, axis, v)
	}
	// write tbody
	for i := 0; i < n; i++ {
		v := valTbody.Index(i)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			continue
		}
		numField := v.NumField() // number of fields in struct

		c := 0
		for j := 0; j < numField; j++ {
			f := v.Field(j)
			t := v.Type().Field(j)
			if !ast.IsExported(t.Name) {
				continue
			}
			tags := ParseTagSetting(t.Tag, "excel")
			// is ignored field
			if _, ok := tags["-"]; ok {
				continue
			}
			c += 1
			var axis string
			axis, err = excelize.CoordinatesToCellName(c, i+2)
			if err != nil {
				err = fmt.Errorf("tbody to cell name error: %s", err)
				return
			}
			_ = xlsx.SetCellValue(sheet, axis, f.Interface())
		}
	}
	// local save
	//xlsx.Path = "./test.xlsx"
	//if err = xlsx.Save(); err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	xlsx.SetActiveSheet(index)
	r, err = xlsx.WriteToBuffer()
	if err != nil {
		err = fmt.Errorf("write to buffer error: %s", err)
		return
	}
	return
}

// ExportSecondaryTitleExcel 导出二级标题.
func ExportSecondaryTitleExcel(firstTitleList []string, secondTitleList [][]string, tbody interface{}) (r *bytes.Buffer, err error) {
	// make sure 'tbody' is a Slice
	valTbody := reflect.ValueOf(tbody)
	if valTbody.Kind() != reflect.Slice {
		err = fmt.Errorf("tbody not slice")
		return
	}
	if len(firstTitleList) != len(secondTitleList) {
		err = fmt.Errorf("title len not equal")
		return
	}
	var (
		// it's a slice, so open up its values
		n     = valTbody.Len()
		sheet = "Sheet1"
	)

	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(sheet)

	col := 0
	for k, title := range firstTitleList {
		if k > 0 {
			col += len(secondTitleList[k-1])
		}
		start, er := excelize.ColumnNumberToName(col + 1)
		if er != nil {
			err = fmt.Errorf("firstTitleList to cell name error: %s", er)
			return
		}
		end, er := excelize.ColumnNumberToName(col + len(secondTitleList[k]))
		if er != nil {
			err = fmt.Errorf("firstTitleList to cell name error: %s", er)
			return
		}
		_ = xlsx.MergeCell(sheet, fmt.Sprintf("%s1", start), fmt.Sprintf("%s1", end))
		var axis string
		axis, err = excelize.CoordinatesToCellName(col+1, 1)
		if err != nil {
			err = fmt.Errorf("firstTitleList to cell name error: %s", err)
			return
		}
		_ = xlsx.SetCellValue(sheet, axis, title)
	}

	// write thead
	col = 0
	for _, list := range secondTitleList {
		for _, subTitle := range list {
			var axis string
			axis, err = excelize.CoordinatesToCellName(col+1, 2)
			if err != nil {
				err = fmt.Errorf("secondTitleList to cell name error: %s", err)
				return
			}
			col++
			_ = xlsx.SetCellValue(sheet, axis, subTitle)
		}
	}

	// write tbody
	for i := 0; i < n; i++ {
		v := valTbody.Index(i)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			continue
		}
		numField := v.NumField() // number of fields in struct

		c := 0
		for j := 0; j < numField; j++ {
			f := v.Field(j)
			t := v.Type().Field(j)
			if !ast.IsExported(t.Name) {
				continue
			}
			tags := ParseTagSetting(t.Tag, "excel")
			// is ignored field
			if _, ok := tags["-"]; ok {
				continue
			}
			c += 1
			var axis string
			axis, err = excelize.CoordinatesToCellName(c, i+3)
			if err != nil {
				err = fmt.Errorf("tbody to cell name error: %s", err)
				return
			}
			_ = xlsx.SetCellValue(sheet, axis, f.Interface())
		}
	}
	xlsx.SetActiveSheet(index)
	r, err = xlsx.WriteToBuffer()
	if err != nil {
		err = fmt.Errorf("write to buffer error: %s", err)
		return
	}
	return
}

// ParseTagSetting get model's field tags.
func ParseTagSetting(tags reflect.StructTag, key ...string) map[string]string {
	setting := map[string]string{}
	for _, v := range key {
		str := tags.Get(v)
		if len(str) == 0 {
			continue
		}
		tags := strings.Split(str, ";")
		for _, value := range tags {
			if len(value) == 0 {
				continue
			}
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}
