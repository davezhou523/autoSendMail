package email

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

// 邮箱配置
const (
	smtpServer  = "smtp.qq.com" // 替换为你的SMTP服务器
	smtpPort    = "587"         // 替换为你的SMTP端口
	senderEmail = "noratf@foxmail.com"
	senderPass  = "qiiqtfkawunibbgb"
)

// knlqvosiwryjbgej
// 收件人列表
var recipients = []string{"davezhou523@gmail.com", "271416962@qq.com"}

// 读取文件内容
func readFileContent(fileName string) (string, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// 发送邮件
func sendEmail(subject, body string) error {
	for _, receiver := range recipients {
		var newReceiver = []string{receiver}
		e := email.NewEmail()
		e.From = fmt.Sprintf("Sender Name <%s>", senderEmail)
		e.To = newReceiver
		e.Subject = subject
		e.Text = []byte(body)

		auth := smtp.PlainAuth("", senderEmail, senderPass, smtpServer)
		var err = e.Send(smtpServer+":"+smtpPort, auth)
		if err != nil {
			return err
		}
	}
	return nil
}

// 定时发送邮件任务
func ScheduleEmail(interval time.Duration, fileName, subject string) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		content, err := readFileContent(fileName)
		if err != nil {
			log.Printf("Failed to read file %s: %v", fileName, err)
			continue
		}
		err = sendEmail(subject, content)
		if err != nil {
			log.Printf("Failed to send email: %v", err)

		} else {
			log.Printf("Email sent successfully with content from %s", fileName)
		}
	}
}
