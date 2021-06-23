package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"reflect"
	"time"
)

type KindType int

const (
	Numbers          KindType = iota // 数字
	lowerCaseLetters                 // 小写字母
	UppercaseLetter                  // 大写字母
	NumberPlusCase                   // 数字加大小写字母
)

// KRand 生成随机字符串
func KRand(size int, kind KindType) []byte {
	iKind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			iKind = KindType(rand.Intn(3))
		}
		scope, base := kinds[iKind][0], kinds[iKind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

// ToMd5 生成md5
func ToMd5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// Assign 批量赋值
func Assign(origin, target interface{}, excludes ...string) {
	valOrigin := reflect.ValueOf(origin).Elem()
	valTarget := reflect.ValueOf(target).Elem()

	for i := 0; i < valOrigin.NumField(); i++ {
		if !valTarget.FieldByName(valOrigin.Type().Field(i).Name).IsValid() {
			continue
		}
		isExclude := false
		for _, col := range excludes {
			if valOrigin.Type().Field(i).Name == col {
				isExclude = true
				break
			}
		}
		if isExclude {
			continue
		}
		tmpOrigin := valOrigin.Field(i)
		tmpTarget := valTarget.FieldByName(valOrigin.Type().Field(i).Name)
		if reflect.TypeOf(tmpOrigin.Interface()) != reflect.TypeOf(tmpTarget.Interface()) {
			continue
		}
		tmpTarget.Set(tmpOrigin)
	}
}
