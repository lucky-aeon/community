package email

import (
	"errors"
	"github.com/jordan-wright/email"
	"net/smtp"
	"time"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/log"
)

var emailPoll *email.Pool

var from string

// 最大重试次数
const maxRetries = 3

// 重试间隔时间
const retryInterval = 2 * time.Second

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
	// 调用带重试的发送方法，默认不重试
	SendWithRetry(to, content, subject, 0)
}

// SendWithRetry 带有重试机制的邮件发送函数
// retries: 重试次数，0表示不重试
func SendWithRetry(to []string, content, subject string, retries int) error {

	// 给自己发，用来测试是否给他人发送成功
	to = append(to, config.GetInstance().EmailConfig.Username)
	if len(to) == 0 {
		return errors.New("收件人列表为空")
	}

	if retries <= 0 {
		retries = 0
	} else if retries > maxRetries {
		retries = maxRetries
	}

	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = subject
	e.Text = []byte(content)

	var err error
	for i := 0; i <= retries; i++ {
		// 第一次尝试或重试
		err = emailPoll.Send(e, 15*time.Second) // 增加超时时间到15秒
		if err == nil {
			return nil // 发送成功
		}

		log.Warnf("发送邮件失败(尝试 %d/%d),接收人: %v,err: %s", i+1, retries+1, to, err.Error())

		// 如果不是最后一次尝试，则等待一段时间后重试
		if i < retries {
			time.Sleep(retryInterval)
		}
	}

	// 批量发送失败后，尝试逐个发送
	return sendIndividually(to, content, subject)
}

// sendIndividually 当批量发送失败时，尝试逐个发送给每个收件人
func sendIndividually(to []string, content, subject string) error {
	log.Infof("批量发送邮件失败，尝试逐个发送给 %d 个收件人", len(to))

	var lastErr error
	successCount := 0

	for _, recipient := range to {
		singleTo := []string{recipient}
		e := email.NewEmail()
		e.From = from
		e.To = singleTo
		e.Subject = subject
		e.Text = []byte(content)

		// 单个发送也尝试重试一次
		var err error
		for i := 0; i <= 1; i++ {
			err = emailPoll.Send(e, 15*time.Second)
			if err == nil {
				successCount++
				break // 发送成功，跳出重试循环
			}

			log.Warnf("单独发送邮件失败(尝试 %d/2),接收人: %s,err: %s", i+1, recipient, err.Error())

			if i < 1 {
				time.Sleep(retryInterval)
			}
		}

		if err != nil {
			lastErr = err
		}
	}

	// 记录单独发送的结果
	if successCount == 0 {
		log.Errorf("所有邮件发送失败，共 %d 个收件人", len(to))
		return lastErr
	} else if successCount < len(to) {
		log.Warnf("部分邮件发送成功: %d/%d", successCount, len(to))
		return errors.New("部分邮件发送失败")
	} else {
		log.Infof("通过单独发送方式，所有邮件发送成功: %d/%d", successCount, len(to))
		return nil
	}
}
