package email

import (
	"fmt"

	"github.com/go-gomail/gomail"
)

// 定义发送邮件内容的结构体
type (
	Conf struct {
		User string
		Name string
		Pwd  string
		Host string
		Port int
	}
	Mail struct {
		To      string
		Subject string
		Msg     string
	}
)

var (
	c = Conf{}
)

// Init 初始化Email配置.
func Init(conf Conf) {
	c = conf
}

func NewMail(to, subject, header, msg, footer string) *Mail {
	return &Mail{To: to, Subject: subject, Msg: GetEmailCurrentMsg(header, msg, footer)}
}

// Send 发送邮件.
func (slf *Mail) Send() (err error) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", c.User, c.Name)
	m.SetHeader("To", m.FormatAddress(slf.To, slf.To))
	m.SetHeader("Subject", slf.Subject)
	m.SetBody("text/html", slf.Msg)
	m.Embed("./mail.png")

	d := gomail.NewDialer(c.Host, c.Port, c.User, c.Pwd)
	if err = d.DialAndSend(m); err != nil {
		return fmt.Errorf("send mail error: %v, mail: %+v", err, slf)
	}
	return
}

// GetEmailCurrentMsg 获取邮箱通用消息模板
func GetEmailCurrentMsg(header, msg, footer string) string {
	return fmt.Sprintf(
		`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Fly</title>
</head>
<body>
<div style="width:638px;margin: 70px auto 0;">
    <div style=" width: 100%%;height: 80px;background-color: #f0f8f6;display: flex;justify-content: left;align-items: center;border-radius: 10px 10px 0 0;">
        <a href="https://github.com/shixiaofeia/fly">
            <img src="cid:mail.png" style=" height: 80px; width: 120px; margin: 0 auto; display: block;border-radius: 10px;" alt=""/>
        </a>
    </div>
    <div style="padding: 30px;min-height: 300px;background: #f1f9f3;font-size: 14px;color: #333;position: relative;border:1px solid #eee;border-radius: 0 0 10px 10px;">
        <p style="margin-bottom:22px;font-size: 16px;">%s</p>
        <p style="position: relative;margin-bottom:16px;">
            %s
        </p>
        <div style="position: absolute;left: 30px;bottom: 30px;color: #999999;font-size: 12px;">%s</div>

    </div>
</div>
</body>
</html>`, header, msg, footer)
}
