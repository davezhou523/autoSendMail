package logic

import (
	"automail/db"
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"strings"
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
	contract, err := l.svcCtx.SearchContact.FindAll(l.ctx)
	if err != nil {
		return
	}
	for _, value := range contract {
		if value.Email == "" {
			continue
		}
		var isReplay uint64 = 1
		task, err := l.svcCtx.EmailTask.FindAll(l.ctx, value.Email, isReplay)
		if err != nil {
			return
		}
		for _, v := range task {
			//是否回复,1:未回复，2：已回复
			if v.IsReplay == 1 {
				continue
			}
		}
		if len(task) == 0 {

		}
	}
	//rows, err := db.DB.Query("SELECT id,title,content,attach_id FROM  email_content ")
	//if err != nil {
	//	fmt.Printf("查询失败: %v\n", err)
	//	return
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var id int
	//	var title string
	//	var content string
	//	var attach_id string
	//	if err := rows.Scan(&id, &title, &content, &attach_id); err != nil {
	//		log.Fatalf("Scan 失败: %v", err)
	//	}
	//
	//	attch, err := getAttch(attach_id)
	//	if err != nil {
	//		return
	//	}
	//	go ScheduleEmail(1*time.Second, content, title, &attch)
	//	//fmt.Printf("email_content: %d, %s\n", id, attach_id)
	//}

	//go email.ScheduleEmail(5*24*time.Hour, "content2.txt", "Content 2")
	//var recipients = []string{"a@gmail.com", "b@gmail.com"}
	//for key, receiver := range recipients {
	//	fmt.Println(key, receiver)
	//}
	// 保持程序运行
	//select {}
	//time.Sleep(10 * time.Second)
	return
}
func getAttch(attach_id string) ([]Attach_struct, error) {
	var ids []int
	err := json.Unmarshal([]byte(attach_id), &ids)
	if err != nil {
		fmt.Printf("attach_id json失败: %v\n", err)
		return nil, err
	}
	// 构建查询语句
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?" // 占位符
		args[i] = id          // 参数
	}
	query := fmt.Sprintf("SELECT  file_name,file_path  FROM  attach where id IN (%s)", strings.Join(placeholders, ", "))
	rows, err := db.DB.Query(query, args...)
	defer rows.Close()
	if err != nil {
		fmt.Printf("attach查询失败: %v\n", err)
		return nil, err
	}
	attachArr := make([]Attach_struct, 0)
	for rows.Next() {
		attach := Attach_struct{}
		if err := rows.Scan(&attach.file_name, &attach.file_path); err != nil {
			log.Fatalf("Scan 失败: %v", err)
		}
		attachArr = append(attachArr, attach)
		fmt.Printf("attach: %s, %s\n", attach.file_name, attach.file_path)
	}
	return attachArr, nil

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
func sendEmail(subject, body string, attach *[]Attach_struct) error {
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
		for _, attach := range *attach {
			m.Embed("." + attach.file_path)
		}
		// 创建并配置邮件拨号器
		d := gomail.NewDialer(smtpServer, smtpPort, senderEmail, senderPass)
		// 发送邮件
		if err := d.DialAndSend(m); err != nil {
			log.Fatalf("send mail fail: %v", err)
			return err
		}
		fmt.Println("send mail finsh")
	}
	return nil
}

// 定时发送邮件任务
func ScheduleEmail(interval time.Duration, content, title string, attach *[]Attach_struct) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		err := sendEmail(title, content, attach)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email sent successfully with content from %s", title)
		}
	}
}
