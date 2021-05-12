package email

import (
	"fmt"
	"github.com/go-gomail/gomail"
)

// 定义发送邮件内容的结构体
type SenEmailMsg struct {
	To    string
	Title string
	Msg   string
}

type Conf struct {
	Website  string
	User     string
	NickName string
	Password string
	Logo     string
}

var (
	c = Conf{}
)

// Init 初始化Email配置
func Init(conf Conf) {
	c = conf
}

// SendEmail 发送邮件
func (s *SenEmailMsg) SendEmail() (err error) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", c.User, c.NickName)
	m.SetHeader("To", m.FormatAddress(s.To, s.To))
	m.SetHeader("Subject", s.Title)
	m.SetBody("text/html", GetEmailCurrentMsg(s.Msg))

	d := gomail.NewDialer("smtpdm.aliyun.com", 80, c.User, c.Password)
	if err = d.DialAndSend(m); err != nil {
		return fmt.Errorf("SendEmail error: %v", err)
	}
	return
}

// GetEmailCurrentMsg 获取邮箱通用消息模板
func GetEmailCurrentMsg(msg string) string {
	return fmt.Sprintf(
		`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>TokenPay</title>
</head>
<body>
    <div style="margin:0 auto;width:638px;margin-top:70px">
        <a href = `+c.Website+` target="_blank">
        <div style=" width: 100%%;height: 80px;background-color: #28292D;display: flex;justify-content: center;align-items: center;">
				<img src=`+c.Logo+` style=" height:41px;"/>
        </div>
	  </a>
             <div style="padding: 30px;min-height: 300px;background: #fff;font-size: 14px;color: #333;position: relative;border:1px solid #eee;">
            <p style="margin-bottom:22px;font-size: 16px;">尊敬的fly开发者，</p>
            <p style="position: relative;left:-5px;margin-bottom:16px;">
               %v
            </p>
            <div style="position: absolute;left: 30px;bottom: 30px;color: #999999;font-size: 12px;">感谢您选择 fly！</div>

        </div>
    </div>
</body>
</html>`, msg)
}
