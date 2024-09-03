package logic

import (
	"automail/autoMail/internal/svc"
	"automail/model"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/bytbox/go-pop3"
	"github.com/jhillyerd/enmime"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

type AutoMailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 邮箱配置
// const (
//
//	smtpServer  = "smtp.qq.com" // 替换为你的SMTP服务器
//	smtpPort    = 587           // 替换为你的SMTP端口
//	senderEmail = "noratf@foxmail.com"
//	senderPass  = "qiiqtfkawunibbgb"
//
// )
const (
	smtpServer  = "smtp.163.com" // 替换为你的SMTP服务器
	smtpPort    = 465            // 替换为你的SMTP端口
	senderEmail = "sunweiglove@163.com"
	senderPass  = "TYKXQAHLUFLVWOFC"
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
	//is_send 是否发送邮件,1:发送，2：不发送
	var isSend uint64 = 1
	//分类,1:手动,2:google
	var category uint64 = 0
	email := "notEmpty"
	var page uint64 = 1
	var pageSize uint64 = 50
	for {
		contract, err := l.svcCtx.SearchContact.FindAll(l.ctx, isSend, category, email, page, pageSize)
		page = page + 1
		if len(contract) == 0 {
			msg := "未查询到需要发送邮件的客户"
			l.Logger.Infof(msg)
			fmt.Println(msg)
			break
		}

		if !errors.Is(err, model.ErrNotFound) && err != nil {
			l.Logger.Error(err)
			break
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
				continue
			}
			if task == nil {
				//查询第一封邮件内容
				fmt.Println("查询第一封邮件内容" + customer.Email)
				emailContent, err := l.svcCtx.EmailContent.FindOneBySort(l.ctx, 1)
				if err != nil {
					l.Logger.Error(err)
					continue
				}
				l.handleSendmail(customer, emailContent)
			} else {
				//查询第下一封邮件内容
				currentEmailContent, err := l.svcCtx.EmailContent.FindOne(l.ctx, task.ContentId)
				//获取下一封要发邮件
				nextSort := currentEmailContent.Sort + 1
				emailContent, err := l.svcCtx.EmailContent.FindOneBySort(l.ctx, nextSort)
				if errors.Is(err, model.ErrNotFound) {
					//is_send 是否发送邮件,1:发送，2：不发送
					customer.IsSend = 2
					err := l.svcCtx.SearchContact.Update(l.ctx, customer)
					if err != nil {
						l.Logger.Error(err)
						continue
					}

					fmt.Printf("%v 所有邮件内容已发送完\n", customer.Email)
					l.Logger.Infof("%v 所有邮件内容已发送完\n", customer.Email)
					continue
				}
				if err != nil {
					l.Logger.Errorf("next emailContent %v", err)
					continue
				}
				l.handleSendmail(customer, emailContent)
			}
		}
		// 添加延迟，避免一次发送太多邮件
		time.Sleep(2 * time.Second)
	}

}

// 邮箱域名转小写
func (l *AutoMailLogic) ConvertEmailDomainLower(customer *model.SearchContact) error {
	parts := strings.Split(customer.Email, "@")
	if len(parts) == 2 {
		parts[1] = strings.ToLower(parts[1]) // 仅将域名部分转为小写
	} else {
		return nil
	}
	customer.Email = strings.Join(parts, "@")
	l.svcCtx.SearchContact.Update(l.ctx, customer)
	return nil
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
	if err != nil {
		return
	}
	go func() {
		err := sendEmail(customer.Email, emailContent.Title, emailContent.Content, attach)
		if err != nil {
			l.Logger.Errorf("sendmail:%v", err)
			return
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
		//fmt.Println("." + attach.FilePath)
		m.Embed("." + attach.FilePath)
	}
	// 创建并配置邮件拨号器
	d := gomail.NewDialer(smtpServer, smtpPort, senderEmail, senderPass)
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("send mail fail: %v", err)
		return err
	}
	fmt.Println(receiver + " send mail finsh")
	return nil
}

func (l *AutoMailLogic) ReceiveEmail() {
	// 设置 POP3 服务器和登录信息
	//	smtpServer  = "smtp.qq.com" // 替换为你的SMTP服务器
	//	smtpPort    = 587           // 替换为你的SMTP端口
	//	senderEmail = "noratf@foxmail.com"
	//	senderPass  = "qiiqtfkawunibbgb"
	//
	pop3Server := "pop.qq.com:995" // 使用POP3的服务器地址和端口
	username := "noratf@foxmail.com"
	password := "qiiqtfkawunibbgb"

	// 建立TLS连接
	conn, err := tls.Dial("tcp", pop3Server, &tls.Config{})
	if err != nil {
		log.Fatal("Failed to connect to POP3 server:", err)
	}
	defer conn.Close()

	// 创建POP3客户端
	client, err := pop3.NewClient(conn)
	if err != nil {
		log.Fatal("Failed to create POP3 client:", err)
	}

	// 用户登录
	if err := client.Auth(username, password); err != nil {
		log.Fatal("Failed to authenticate:", err)
	}

	// 获取邮箱状态
	count, size, err := client.Stat()
	if err != nil {
		log.Fatal("Failed to get mailbox status:", err)
	}
	fmt.Printf("You have %d messages, total size is %d bytes.\n", count, size)

	// POP3协议中，邮件编号是按时间顺序排列的，编号越大，邮件越新。因此，你可以从最大的编号开始遍历，直到找到符合条件的邮件。
	for i := count; i > 1080; i-- {
		//// 获取邮件头部信息
		//header, err := client.Top(i, 0)
		//
		//if err != nil {
		//	log.Printf("Failed to retrieve message %d: %v\n", i, err)
		//	continue
		//}
		//fmt.Printf("Message %d Header:\n%s\n", i, header)

		// 获取完整邮件内容
		msg, err := client.Retr(i)
		if err != nil {
			log.Printf("Failed to retrieve message %d: %v\n", i, err)
			continue
		}
		// 使用 enmime 解析邮件内容
		reader := strings.NewReader(msg)
		env, err := enmime.ReadEnvelope(reader)
		if err != nil {
			log.Printf("Failed to parse message %d: %v\n", i, err)
			continue
		}

		// 输出邮件主题和正文
		//for _, v := range env.GetHeaderKeys() {
		//	//fmt.Println(env.GetHeader(v))
		//	//获取头部
		//	//X-Qq-Action
		//	//X-Qq-Style
		//	//From
		//	//Message-Id
		//	//X-Qq-Mime
		//	//X-Mailer
		//	//X-Qq-Mailer
		//	//Mime-Version
		//	//Content-Transfer-Encoding
		//	//Date
		//	//X-Qq-Mid
		//	//X-Qq-Inner-Pending
		//	//To
		//	//Subject
		//	//Content-Type
		//}
		//fmt.Println(env.GetHeader("From"), env.GetHeader("Date"))
		from := env.GetHeader("From")
		fmt.Println(from)
		if strings.Contains(from, "PostMaster@qq.com") {
			//fmt.Printf("Message %d Text Body: %s\n", i, env.Text)
			// 电子邮件地址正则表达式
			emailRegex := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`

			// 编译正则表达式
			re := regexp.MustCompile(emailRegex)
			emails := re.FindString(env.Text)
			searchContactList, _ := l.svcCtx.SearchContact.FindAll(l.ctx, 0, 0, emails, 1, 1)
			fmt.Println(emails, searchContactList, len(searchContactList))
			if len(searchContactList) > 0 {
				fmt.Println("系统存在:" + emails)
				for _, searchContact := range searchContactList {
					err := l.svcCtx.SearchContact.Delete(l.ctx, searchContact.Id)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

			}
			// 你可以选择删除邮件
			err = client.Dele(i)
			if err != nil {
				log.Fatal("Failed to delete message:", err)
			}
		}

		//fmt.Printf("Message %d Text Body: %s\n", i, env.Text)
		//fmt.Printf("Message %d HTML Body: %s\n", i, env.HTML)

	}

	fmt.Println("Done!")
}
