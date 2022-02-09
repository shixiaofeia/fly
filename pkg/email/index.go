package email

import (
	"fmt"

	"github.com/go-gomail/gomail"
)

// 定义发送邮件内容的结构体
type (
	Conf struct {
		Website  string
		User     string
		NickName string
		Password string
		Logo     string
	}
	SenEmailMsg struct {
		To    string
		Title string
		Msg   string
	}
)

var (
	c = Conf{}
)

// Init 初始化Email配置.
func Init(conf Conf) {
	c = conf
}

// SendEmail 发送邮件.
func (s *SenEmailMsg) SendEmail() (err error) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", c.User, c.NickName)
	m.SetHeader("To", m.FormatAddress(s.To, s.To))
	m.SetHeader("Subject", s.Title)
	m.SetBody("text/html", s.Msg)

	d := gomail.NewDialer("smtpdm.aliyun.com", 80, c.User, c.Password)
	if err = d.DialAndSend(m); err != nil {
		return fmt.Errorf("SendEmail error: %v", err)
	}
	return
}
