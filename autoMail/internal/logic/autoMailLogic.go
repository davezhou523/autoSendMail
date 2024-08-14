package logic

import (
	"automail/model"
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"time"

	"automail/autoMail/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AutoMailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}
type Attach_struct struct {
	file_name string
	file_path string
}

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

func NewAutoMailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AutoMailLogic {
	return &AutoMailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AutoMailLogic) AutoMail() {
	//isReplay是否回复,1:未回复，2：已回复
	var isReplay uint64 = 1
	contract, err := l.svcCtx.SearchContact.FindAll(l.ctx, isReplay)
	if err != nil {
		return
	}
	for _, customer := range contract {
		if customer.Email == "" {
			continue
		}

		task, err := l.svcCtx.EmailTask.FindAll(l.ctx, customer.Email)
		if err != nil {
			return
		}
		if len(task) == 0 {
			//查询第一封邮件内容
			emailContent, err := l.svcCtx.EmailContent.FindOneBySort(l.ctx, 1)
			if err != nil {
				return
			}
			attach, err := l.getAttach(emailContent.AttachId)
			fmt.Println(attach)
			if err != nil {
				return
			}
			go func() {
				err := sendEmail(customer.Email, emailContent.Title, emailContent.Content, attach)
				if err != nil {
					fmt.Println(err)
				}
				emailTask := new(model.EmailTask)
				emailTask.Email = customer.Email
				emailTask.ContentId = emailContent.Id
				emailTask.SendTime = time.Now().Unix()
				emailTask.CreateTime = time.Now().Format("2006-01-02 15:04:05")
				et, err := l.svcCtx.EmailTask.Insert(l.ctx, emailTask)
				id, err := et.LastInsertId()
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("LastInsertId:%d", id)
				if err != nil {
					return
				}
			}()
		} else {
			//for _, v := range task {
			//
			//}
		}

	}

	return
}
func (l *AutoMailLogic) getAttach(attach_id string) ([]*model.Attach, error) {
	attach, err := l.svcCtx.Attach.FindAll(l.ctx, attach_id)
	if err != nil {
		return nil, err
	}
	return attach, nil
	//attachArr := make([]Attach_struct, 0)
	//for rows.Next() {
	//	attach := Attach_struct{}
	//	if err := rows.Scan(&attach.file_name, &attach.file_path); err != nil {
	//		log.Fatalf("Scan 失败: %v", err)
	//	}
	//	attachArr = append(attachArr, attach)
	//	fmt.Printf("attach: %s, %s\n", attach.file_name, attach.file_path)
	//}
	//return attachArr, nil

}
func saveEmailTask() {

}

// 读取文件内容
func readFileContent(fileName string) (string, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// 发送邮件
func sendEmail(receiver, subject, body string, attach []*model.Attach) error {
	// 创建新的消息
	m := gomail.NewMessage()
	// 设置邮件头
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", subject)

	// 设置邮件主体内容（HTML格式）
	m.SetBody("text/html", body)

	// 添加图片（内嵌图片）
	for _, attach := range attach {
		fmt.Println("." + attach.FilePath)
		m.Embed("." + attach.FilePath)
	}
	// 创建并配置邮件拨号器
	d := gomail.NewDialer(smtpServer, smtpPort, senderEmail, senderPass)
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("send mail fail: %v", err)
		return err
	}
	fmt.Println("send mail finsh")

	return nil
}

// 定时发送邮件任务
func ScheduleEmail(interval time.Duration, content, title string, attach []Attach_struct) {
	//ticker := time.NewTicker(interval)
	//for range ticker.C {
	//	err := sendEmail("", title, content, attach)
	//	if err != nil {
	//		log.Printf("Failed to send email: %v", err)
	//	} else {
	//		log.Printf("Email sent successfully with content from %s", title)
	//	}
	//}
}
