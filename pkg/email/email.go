package email

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"xhyovo.cn/community/pkg/log"
)

var (
	smtpHost     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
)

// Init 初始化邮件发送配置
func Init(username, password, host string) {
	smtpHost = host
	smtpUsername = username
	smtpPassword = password
	fromEmail = username
	log.Infof("邮件服务初始化完成，SMTP服务器: %s, 用户名: %s", smtpHost, smtpUsername)
}

// Send 发送邮件
func Send(to []string, content, subject string) error {
	// 检查收件人列表
	if len(to) == 0 {
		return fmt.Errorf("收件人列表为空")
	}

	// 构建邮件内容
	date := fmt.Sprintf("%s", time.Now().Format(time.RFC1123Z))
	toAddress := strings.Join(to, ";")

	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", fromEmail, smtpUsername)
	header["To"] = toAddress
	header["Subject"] = subject
	header["Date"] = date
	header["Content-Type"] = "text/plain; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + content

	// 先尝试批量发送
	err := smtp.SendMail(
		smtpHost,
		smtp.PlainAuth("", smtpUsername, smtpPassword, strings.Split(smtpHost, ":")[0]),
		smtpUsername,
		to,
		[]byte(message),
	)

	if err == nil {
		log.Infof("邮件批量发送成功，收件人: %v", to)
		return nil
	}

	// 批量发送失败，记录错误
	log.Warnf("邮件批量发送失败，错误: %s，尝试单独发送", err.Error())

	// 批量发送失败后，尝试逐个发送
	return sendIndividually(to, content, subject)
}

// sendIndividually 当批量发送失败时，尝试逐个发送给每个收件人
func sendIndividually(to []string, content, subject string) error {
	log.Infof("开始逐个发送邮件给 %d 个收件人", len(to))

	var lastErr error
	successCount := 0

	for _, recipient := range to {
		singleTo := []string{recipient}

		// 构建邮件内容
		date := fmt.Sprintf("%s", time.Now().Format(time.RFC1123Z))

		header := make(map[string]string)
		header["From"] = fmt.Sprintf("%s <%s>", fromEmail, smtpUsername)
		header["To"] = recipient
		header["Subject"] = subject
		header["Date"] = date
		header["Content-Type"] = "text/plain; charset=UTF-8"

		message := ""
		for k, v := range header {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + content

		// 尝试单独发送
		err := smtp.SendMail(
			smtpHost,
			smtp.PlainAuth("", smtpUsername, smtpPassword, strings.Split(smtpHost, ":")[0]),
			smtpUsername,
			singleTo,
			[]byte(message),
		)

		if err == nil {
			successCount++
			log.Infof("单独发送邮件成功，接收人: %s", recipient)
		} else {
			lastErr = err
			log.Warnf("单独发送邮件失败，接收人: %s，错误: %s", recipient, err.Error())
		}

		// 单个发送之间添加短暂延迟，避免频率限制
		time.Sleep(500 * time.Millisecond)
	}

	// 记录单独发送的结果
	if successCount == 0 {
		log.Errorf("所有邮件发送失败，共 %d 个收件人", len(to))
		return lastErr
	} else if successCount < len(to) {
		log.Warnf("部分邮件发送成功: %d/%d", successCount, len(to))
		return fmt.Errorf("部分邮件发送失败")
	} else {
		log.Infof("通过单独发送方式，所有邮件发送成功: %d/%d", successCount, len(to))
		return nil
	}
}
