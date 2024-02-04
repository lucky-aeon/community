package email

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"time"
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
		log.Fatalln("email connect fail", err.Error())
		return
	}
	emailPoll = p
	from = username
}

func Send(to []string, content, subject string) {
	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = subject
	e.Text = []byte(content)
	err := emailPoll.Send(e, 10*time.Second)
	if err != nil {
		return
	}
}
