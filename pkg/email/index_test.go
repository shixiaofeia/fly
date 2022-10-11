package email

import (
	"os"
	"strings"
	"testing"
)

func TestMail_Send(t *testing.T) {
	conf := Conf{
		User: "",
		Name: "Fly",
		Pwd:  "",
		Host: "smtpdm.aliyun.com",
		Port: 80,
	}
	Init(conf)

	if err := NewMail("", "Fly通知", "尊敬的Fly用户,", "这是一封测试邮件", "感谢您的选择!").Send(); err != nil {
		t.Errorf("send email err: %v", err)
		return
	}

	t.Log("send email success")
}

func TestGetEmailCurrentMsg(t *testing.T) {
	msg := GetEmailCurrentMsg("尊敬的用户", "您好", "感谢您的选择")
	msg = strings.Replace(msg, "cid:mail.png", "mail.png", 1)
	file, err := os.OpenFile("./mail.html", os.O_WRONLY, 0744)
	if err != nil {
		t.Errorf("open mail html err: %v", err)
		return
	}
	if _, err = file.WriteString(msg); err != nil {
		t.Errorf("write html err: %v", err)
	}
}
