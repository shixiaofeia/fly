package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type InviteCode struct {
	base    string // 进制的包含字符, string类型
	decimal uint32 // 进制长度
	pad     string // 补位字符,若生成的code小于最小长度,则补位+随机字符, 补位字符不能在进制字符中
	len     int    // code最小长度
}

// NewInviteCode 邀请码
func NewInviteCode() *InviteCode {
	return &InviteCode{
		base:    "W7PD9CZS2VE8N6LH1XMJ4TUK3ARG5YQI",
		decimal: 32,
		pad:     "B",
		len:     8,
	}
}

// IdToCode id转code
func (c *InviteCode) IdToCode(id uint32) string {
	mod := uint32(0)
	res := ""
	for id != 0 {
		mod = id % c.decimal
		id = id / c.decimal
		res += string(c.base[mod])
	}
	resLen := len(res)
	if resLen < c.len {
		res += c.pad
		for i := 0; i < c.len-resLen-1; i++ {
			rand.Seed(time.Now().UnixNano())
			res += string(c.base[rand.Intn(int(c.decimal))])
		}
	}
	return res
}

// CodeToId code转id
func (c *InviteCode) CodeToId(code string) uint32 {
	res := uint32(0)
	lenCode := len(code)

	baseArr := []byte(c.base)     // 字符串进制转换为byte数组
	baseRev := make(map[byte]int) // 进制数据键值转换为map
	for k, v := range baseArr {
		baseRev[v] = k
	}

	// 查找补位字符的位置
	isPad := strings.Index(code, c.pad)
	if isPad != -1 {
		lenCode = isPad
	}

	r := 0
	for i := 0; i < lenCode; i++ {
		// 补充字符直接跳过
		if string(code[i]) == c.pad {
			continue
		}
		index := baseRev[code[i]]
		b := uint32(1)
		for j := 0; j < r; j++ {
			b *= c.decimal
		}
		res += uint32(index) * b
		r++
	}
	return res
}

// initCheck 初始化检查
func (c *InviteCode) initCheck() (bool, error) {
	lenBase := len(c.base)
	// 检查进制字符
	if c.base == "" {
		return false, fmt.Errorf("base string is nil or empty")
	}
	// 检查长度是否符合
	if uint32(lenBase) != c.decimal {
		return false, fmt.Errorf("base length and len not match")
	}
	return true, nil
}
