package structf

import (
	"encoding/json"
	"reflect"
)

// Assign 赋值结构体并忽略指定字段.
// 适用于结构体转换, 想要使用更全面的转换建议使用 https://goframe.org/pages/viewpage.action?pageId=1114677.
func Assign(src, dst interface{}, excludes ...string) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		value := srcVal.Field(i)
		name := srcVal.Type().Field(i).Name

		isExclude := false
		for _, col := range excludes {
			if col == name {
				isExclude = true
				break
			}
		}
		if isExclude {
			continue
		}

		dstValue := dstVal.FieldByName(name)
		if !dstValue.IsValid() {
			continue
		}
		dstValue.Set(value)
	}
}

// AssignUnmarshal json赋值.
func AssignUnmarshal(src, dst interface{}) {
	b, _ := json.Marshal(src)
	_ = json.Unmarshal(b, dst)
}
