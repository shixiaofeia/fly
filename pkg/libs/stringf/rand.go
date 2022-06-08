package stringf

import (
	"math/rand"
	"time"
)

type KindType int

const (
	Numbers          KindType = iota // 数字
	LowerCaseLetters                 // 小写字母
	UppercaseLetter                  // 大写字母
	NumberPlusCase                   // 数字加大小写字母
)

// RandStr 生成随机字符串.
func RandStr(size int, kind KindType) []byte {
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
