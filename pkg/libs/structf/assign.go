package structf

import "reflect"

// Assign 赋值结构体并忽略指定字段.
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
