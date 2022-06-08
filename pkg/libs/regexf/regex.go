package regexf

import (
	"github.com/dlclark/regexp2"
	"regexp"
)

type Regex struct{}

var _regex = new(Regex)

func NewRegex() *Regex {
	return _regex
}

// Pwd 密码大小写字母+数字+特殊字符, 最少包含三类.
func (*Regex) Pwd(pwd string) bool {
	expr := `^(?![a-zA-Z]+$)(?![A-Z0-9]+$)(?![A-Z\W_]+$)(?![a-z0-9]+$)(?![a-z\W_]+$)(?![0-9\W_]+$)[a-zA-Z0-9\W_]{8,20}$`
	reg, _ := regexp2.Compile(expr, 0)
	m, _ := reg.FindStringMatch(pwd)
	if m != nil {
		return true
	}
	return false
}

// Email email.
func (*Regex) Email(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
