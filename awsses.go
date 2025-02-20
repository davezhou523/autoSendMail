package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

func main() {
	// 创建 AWS SES 会话
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("southeast-2"), // 选择你申请的 AWS SES 区域
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// 初始化 SES 客户端
	svc := ses.New(sess)

	// 配置邮件内容
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String("janiehuang@tenfangmt.com"), // 收件人邮箱
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String("Hello, this is a test email from Amazon SES!"),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Test Email from SES"),
			},
		},
		Source: aws.String("yourname@yourcompany.com"), // 你的发件邮箱（必须已验证）
	}

	// 发送邮件
	result, err := svc.SendEmail(input)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully:", result)
}
