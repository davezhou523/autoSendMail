package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"time"
)

// 邮箱配置
const (
	smtpServer  = "smtp.qq.com" // 替换为你的SMTP服务器
	smtpPort    = 587           // 替换为你的SMTP端口
	senderEmail = "noratf@foxmail.com"
	senderPass  = "qiiqtfkawunibbgb"
)

// knlqvosiwryjbgej
// 收件人列表
var recipients = []string{"davezhou523@gmail.com", "271416962@qq.com", "731847483@qq.com"}

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
		// 创建新的消息
		m := gomail.NewMessage()
		// 设置邮件头
		m.SetHeader("From", senderEmail)
		m.SetHeader("To", receiver)
		m.SetHeader("Subject", subject)

		// 设置邮件主体内容（HTML格式）
		m.SetBody("text/html", body)

		// 添加图片（内嵌图片）
		//m.Embed("./static/content2-1.png")
		//m.Embed("./static/content2-2.png")
		//m.Embed("./static/content2-3.png")
		//m.Embed("./static/content2-4.png")
		m.Embed("./static/content4-1.png")
		m.Embed("./static/content4-2.png")
		m.Embed("./static/content4-3.png")
		// 创建并配置邮件拨号器
		d := gomail.NewDialer(smtpServer, smtpPort, senderEmail, senderPass)

		// 发送邮件
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
		fmt.Println("send mail finsh")
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
