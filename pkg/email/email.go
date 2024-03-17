package email

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"time"
	"xhyovo.cn/community/pkg/log"
)

var emailPoll *email.Pool

var from string

func Init(address, username, password, host string, pollCount int) {
	p, err := email.NewPool(
		address,
		pollCount,
		smtp.PlainAuth("", username, password, host),
	)
	if err != nil {
		log.Errorf("初始化 email 失败,err: %s", err.Error())
		panic(err.Error())
		return
	}
	emailPoll = p
	from = username
}

func Send(to []string, content, subject string) {
	if len(to) == 0 {
		return
	}
	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = subject
	e.Text = []byte(content)
	err := emailPoll.Send(e, 10*time.Second)
	if err != nil {
		log.Warnf("发送邮箱失败,接收人: %v,err: %s", to, err.Error())
		return
	}
}
