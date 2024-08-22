package logic

import (
	"automail/autoMail/internal/svc"
	"automail/model"
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
)

type AutoMailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
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
	//分类,1:手动,2:google
	var category uint64 = 1
	contract, err := l.svcCtx.SearchContact.FindAll(l.ctx, isReplay, category)
	if !errors.Is(err, model.ErrNotFound) && err != nil {
		l.Logger.Error(err)
		return
	}
	for _, customer := range contract {
		if customer.Email == "" {
			continue
		}
		fmt.Printf("customer email:%v\n", customer.Email)
		//通过email查最新发邮件任务的记录
		task, err := l.svcCtx.EmailTask.FindOneBySort(l.ctx, 0, customer.Email)
		if !errors.Is(err, model.ErrNotFound) && err != nil {
			l.Logger.Error(err)
			return
		}
		if task == nil {
			//查询第一封邮件内容
			fmt.Println("查询第一封邮件内容" + customer.Email)
			emailContent, err := l.svcCtx.EmailContent.FindOneBySort(l.ctx, 1)
			if err != nil {
				l.Logger.Error(err)
				return
			}
			l.handleSendmail(customer, emailContent)
		} else {
			//查询第下一封邮件内容
			currentEmailContent, err := l.svcCtx.EmailContent.FindOne(l.ctx, task.ContentId)
			//获取下一封要发邮件
			nextSort := currentEmailContent.Sort + 1
			emailContent, err := l.svcCtx.EmailContent.FindOneBySort(l.ctx, nextSort)
			if errors.Is(err, model.ErrNotFound) {
				l.Logger.Errorf("%v 所有邮件内容已发送完\n", customer.Email)
				continue
			}
			if err != nil {
				l.Logger.Errorf("next emailContent %v", err)
				return
			}
			l.handleSendmail(customer, emailContent)
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
}
func (l *AutoMailLogic) handleSendmail(customer *model.SearchContact, emailContent *model.EmailContent) {
	attach, err := l.getAttach(emailContent.AttachId)
	fmt.Println(attach)
	if err != nil {
		return
	}
	go func() {
		err := sendEmail(customer.Email, emailContent.Title, emailContent.Content, attach)
		if err != nil {
			l.Logger.Errorf("sendmail:%v", err)
		}
		id, err := NewEmailTaskLogic(l.ctx, l.svcCtx).saveEmailTask(customer, emailContent)
		if err != nil {
			l.Logger.Errorf("saveEmailTask:%v", err)
			return
		}

		fmt.Printf("LastInsertId:%d\n", id)
		if err != nil {
			return
		}
	}()
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
	fmt.Println(" send mail finsh")

	return nil
}
