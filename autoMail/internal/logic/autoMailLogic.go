package logic

import (
	"automail/autoMail/internal/svc"
	"automail/model"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/bytbox/go-pop3"
	"github.com/jhillyerd/enmime"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/semaphore"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

type AutoMailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 收件人列表
var recipients = []string{}

var sem = semaphore.NewWeighted(10) // 最多允许 10 个协程同时发送邮件
var wg sync.WaitGroup

func NewAutoMailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AutoMailLogic {
	return &AutoMailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AutoMailLogic) worker(wg *sync.WaitGroup, customerTasks chan *model.SearchContact, providers []*model.EmailProviders, emailContent *model.EmailContent) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			l.Logger.Errorf("worker recover from panic:%v", r)
		}
	}()
	providerIndex := 0
	for customer := range customerTasks {
		if len(providers) == 0 {
			fmt.Println("没有可用的邮件服务商，任务暂停")
			break
		}
		provider := providers[providerIndex]
		fmt.Printf("邮件服务商 Email:%v,客户 Email:%v\n", provider.Username, customer.Email)
		l.handleSendmail(provider, customer, emailContent)
		providerIndex = (providerIndex + 1) % len(providers) // 轮询选择 SMTP 账号
		fmt.Printf("协程数:%v\n", runtime.NumGoroutine())

	}
}

func (l *AutoMailLogic) AutoMail() {
	//is_send 是否发送邮件,1:发送，2：不发送
	//分类,1:手动,2:google
	var category int64 = 0
	var company_id int64 = 1
	var user_id int64 = 1
	email := ""
	var page uint64 = 1
	var pageSize uint64 = 10
	//var sort uint64 = 5
	create_time := "2025-02-12"
	var contentId uint64 = 7
	emailContent, _ := l.svcCtx.EmailContent.FindOne(l.ctx, contentId)
	if emailContent == nil {
		l.Logger.Errorf("邮件模板内容不存在,id：%v\n", contentId)
		return
	}
	for {
		contacts, err := l.svcCtx.SearchContact.FindAll(l.ctx, user_id, company_id, category, 0, email, create_time, page, pageSize, contentId)
		page = page + 1
		if len(contacts) == 0 {
			msg := "未查询到需要发送邮件的客户"
			l.Logger.Infof(msg)
			fmt.Println(msg)
			break
		}

		if !errors.Is(err, model.ErrNotFound) && err != nil {
			l.Logger.Error(err)
			break
		}
		// 检查是否需要重置限额
		_, _ = l.svcCtx.EmailProviders.ResetDailyCount()
		providers, err := NewEmailProvidersLogic(l.ctx, l.svcCtx).getProvidersList(user_id, company_id)
		if err != nil {
			l.Logger.Error(err.Error())
			return
		}
		var wg sync.WaitGroup
		workerCount := len(providers) // 3 个 worker 并发处理
		fmt.Printf("workerCount:%v", workerCount)
		// 创建任务队列
		taskChan := make(chan *model.SearchContact, len(contacts))
		defer func() {
			if r := recover(); r != nil {
				l.Logger.Errorf("recover from panic:%v", r)
			}
		}()
		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go l.worker(&wg, taskChan, providers, emailContent)
		}

		for _, customer := range contacts {
			if customer.Email == "" {
				continue
			}
			taskChan <- customer

			//fmt.Printf("customer email:%v\n", customer.Email)

			//通过email查最新发邮件任务的记录
			//task, err := l.svcCtx.EmailTask.FindOneBySort(l.ctx, 0, customer.Email)
			//if !errors.Is(err, model.ErrNotFound) && err != nil {
			//	l.Logger.Error(err)
			//	continue
			//}
			//if task == nil {
			//	//查询第一封邮件内容
			//	fmt.Println("查询第一封邮件内容" + customer.Email)
			//	emailContent, err := l.svcCtx.EmailContent.FindOneBySort(l.ctx, sort)
			//	if err != nil {
			//		l.Logger.Error("未查询到邮件模板" + err.Error())
			//		continue
			//	}
			//	l.handleSendmail(customer, emailContent)
			//} else {
			//查询第下一封邮件内容
			//currentEmailContent, err := l.svcCtx.EmailContent.FindOne(l.ctx, task.ContentId)
			////获取下一封要发邮件
			//nextSort := currentEmailContent.Sort + 1
			//emailContent, err := l.svcCtx.EmailContent.FindOneBySort(l.ctx, nextSort)
			//if errors.Is(err, model.ErrNotFound) {
			//	//is_send 是否发送邮件,1:发送，2：不发送
			//	customer.IsSend = 2
			//	err := l.svcCtx.SearchContact.Update(l.ctx, customer)
			//	if err != nil {
			//		l.Logger.Error(err)
			//		continue
			//	}
			//
			//	fmt.Printf("%v 所有邮件内容已发送完\n", customer.Email)
			//	l.Logger.Infof("%v 所有邮件内容已发送完\n", customer.Email)
			//	continue
			//}
			//if err != nil {
			//	l.Logger.Errorf("next emailContent %v", err)
			//	continue
			//}
			//l.handleSendmail(customer, emailContent)
			//}

		}
		close(taskChan)
		wg.Wait()

		fmt.Println("📨 所有邮件任务完成")
	}

}

// 账户类型	每秒限制	每分钟限制	每小时限制	每日限制
// 个人 Gmail 账户	约 1-2 封	约 60-100 封	约 100 封	500 封
// GoogleWorkspace约 1-3 封	约 60-150 封	约 200 封	2000 封
func (l *AutoMailLogic) CustomizeSend() {
	//is_send 是否发送邮件,1:发送，2：不发送
	//分类,1:手动,2:google
	var category int64 = 1
	var company_id int64 = 1
	var user_id int64 = 1
	//email := "notEmpty"
	//email := "zhouzeng8709@163.com"
	email := "271416962@qq.com"
	//email := "janiehcn@outlook.com"
	var promotionContentId uint64 = 7 //推广内容id
	var page uint64 = 1
	var pageSize uint64 = 100

	var id uint64 = 0
	for {
		contract, err := l.svcCtx.SearchContact.FindAll(l.ctx, user_id, company_id, category, id, email, "", page, pageSize, promotionContentId)
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
			fmt.Println(customer.Email)
			if customer.Email == "" {
				continue
			}
			fmt.Printf("customer email:%v\n", customer.Email)
			//查询第下一封邮件内容
			//currentEmailContent, err := l.svcCtx.EmailContent.FindOne(l.ctx, promotionContentId)

			if err != nil {
				fmt.Printf("emailContent %v", err)
				l.Logger.Errorf("emailContent %v", err)
				continue
			}
			//l.handleSendmail(customer, currentEmailContent)
		}
		// 添加延迟，避免一次发送太多邮件
		time.Sleep(1 * time.Second)
	}

}

// 邮箱域名转小写
func (l *AutoMailLogic) ConvertEmailDomainLower() error {
	//is_send 是否发送邮件,1:发送，2：不发送
	//分类,1:手动,2:google
	var category int64 = 1
	var company_id int64 = 1
	var user_id int64 = 1
	email := "notEmpty"
	var page uint64 = 1
	var pageSize uint64 = 1000
	var create_time string = "2024-09-14"
	for {
		contract, err := l.svcCtx.SearchContact.FindAll(l.ctx, user_id, company_id, category, 0, email, create_time, page, pageSize, 0)
		page = page + 1
		fmt.Printf("page:%v\n", page)
		if len(contract) == 0 {
			msg := "未查询到需要发送邮件的客户"
			l.Logger.Infof(msg)
			fmt.Println(msg)
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, customer := range contract {

			parts := strings.Split(customer.Email, "@")
			if len(parts) == 2 {
				parts[1] = strings.ToLower(parts[1]) // 仅将域名部分转为小写
			} else {
				return nil
			}
			customer.Email = strings.Join(parts, "@")
			//println(customer.Email)
			err := l.svcCtx.SearchContact.Update(l.ctx, customer)
			if err != nil {
				println(err)
				return err
			}
		}
	}

	return nil
}
func (l *AutoMailLogic) getAttach(attach_id string) ([]*model.Attach, error) {
	attach, err := l.svcCtx.Attach.FindAll(l.ctx, attach_id)
	if err != nil {
		return nil, err
	}
	return attach, nil
}
func (l *AutoMailLogic) validatEmail(customer *model.SearchContact) error {
	// 创建邮箱验证器实例
	ev := emailverifier.NewVerifier()
	// 使用邮箱验证器验证邮箱
	_, err := ev.Verify(customer.Email)
	return err
}

func (l *AutoMailLogic) handleSendmail(provider *model.EmailProviders, customer *model.SearchContact, emailContent *model.EmailContent) {
	vaildateRes := l.validatEmail(customer)
	if vaildateRes != nil {
		l.UpdateReturnByEmail(customer.Email, vaildateRes.Error())
		l.Logger.Errorf("Email is invalid or does not exist:%v", vaildateRes.Error())
		return
	}
	attach, err := l.getAttach(emailContent.AttachId)
	if err != nil {
		return
	}
	time.Sleep(2 * time.Second)
	//
	//defer func() {
	//	if r := recover(); r != nil {
	//		l.Logger.Errorf("recover from panic:%v", r)
	//	}
	//}()

	//重试几次发送
	err = l.sendEmailWithRetry(provider, customer, emailContent, attach, 1)
	//增加发送邮件计数
	//_, _ = l.svcCtx.EmailProviders.IncrementSent(l.ctx, provider.Id)
	if err != nil {
		return
	}
	l.svcCtx.SqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		//增加发送邮件计数
		emailProvider := l.svcCtx.EmailProviders.WithSession(session)
		emailProvider.IncrementSent(l.ctx, provider.Id)

		id, err := NewEmailTaskLogic(l.ctx, l.svcCtx).saveEmailTaskWithSession(session, customer, emailContent)
		fmt.Printf("EmailTask LastInsertId:%d\n", id)
		return err
	})

	//wg.Add(1)
	//go func(customer *model.SearchContact, emailContent *model.EmailContent, attach []*model.Attach, provider *model.EmailProviders) {
	//	defer wg.Done()
	//
	//	defer func() {
	//		if r := recover(); r != nil {
	//			l.Logger.Errorf("recover from panic:%v", r)
	//		}
	//	}()
	//	// 限制并发数量
	//	err := sem.Acquire(l.ctx, 1)
	//	if err != nil {
	//		l.Logger.Errorf("sem.Acquire:%v", err)
	//		return
	//	}
	//	defer sem.Release(1)
	//
	//	//重试几次发送
	//	err = l.sendEmailWithRetry(provider, customer, emailContent, attach, 1)
	//	if err != nil {
	//		return
	//	}
	//	//增加发送邮件发送计数
	//	_, _ = l.svcCtx.EmailProviders.IncrementSent(l.ctx, provider.Id)
	//
	//	id, err := NewEmailTaskLogic(l.ctx, l.svcCtx).saveEmailTask(customer, emailContent)
	//	if err != nil {
	//		l.Logger.Errorf("saveEmailTask:%v", err)
	//		return
	//	}
	//	fmt.Printf("LastInsertId:%d\n", id)
	//	if err != nil {
	//		return
	//	}
	//}(customer, emailContent, attach, provider)
	//wg.Wait()

}

// 重试几次发送
func (l *AutoMailLogic) sendEmailWithRetry(emailProviders *model.EmailProviders, customer *model.SearchContact, emailContent *model.EmailContent, attach []*model.Attach, retries int) error {
	var err error
	for i := 0; i < retries; i++ {
		// 发送邮件逻辑
		//err = l.SendEmail(emailProviders, customer, emailContent, attach)
		if err == nil {
			// 每次发送后增加一个随机延迟，防止频率过高
			time.Sleep(time.Second * time.Duration(rand.Intn(2)+1))
			return nil
		}
		time.Sleep(time.Second * 10) // 等待 2 秒再重试
	}
	fmt.Printf("sendEmailWithRetry err:%v\n", err.Error())
	return err

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
func (l *AutoMailLogic) SendEmail(emailProviders *model.EmailProviders, customer *model.SearchContact, emailContent *model.EmailContent, attach []*model.Attach) error {
	smtpServer := emailProviders.SmtpServer
	smtpPort := emailProviders.SmtpPort
	senderEmail := emailProviders.Username
	senderPass := emailProviders.Password
	unsubscribe := l.svcCtx.Config.Unsubscribe
	replyTo := l.svcCtx.Config.ReplyTo
	//receiver := customer.Email
	receiver := "271416962@qq.com"
	unsubscribeAPI := l.svcCtx.Config.UnsubscribeAPI
	token := "abcdef"
	// 创建新的消息
	m := gomail.NewMessage()
	// 设置邮件头
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", emailContent.Title)
	m.SetHeader("List-Unsubscribe", fmt.Sprintf("<mailto:%v>", unsubscribe))
	m.SetHeader("Subject", replyTo)
	fmt.Println(unsubscribe, replyTo)
	firtname := customer.FirstName
	clientCompany := customer.Company
	unsubscribeUrl := fmt.Sprintf("%s?email=%s&token=%s", unsubscribeAPI, receiver, token)
	mailContent := fmt.Sprintf(emailContent.Content, firtname, clientCompany, unsubscribeUrl)
	// 设置邮件主体内容（HTML格式）
	m.SetBody("text/html", mailContent)

	// 添加图片（内嵌图片）
	for _, attach := range attach {
		//fmt.Println("." + attach.FilePath)
		m.Embed("." + attach.FilePath)
	}
	// 创建并配置邮件拨号器

	d := gomail.NewDialer(smtpServer, int(smtpPort), senderEmail, senderPass)
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		//550 User is over flow 错误通常表示收件人的邮箱已满，无法接收更多邮件。

		if err.Error() == "550 User is over flow" {
			//系统退回0:未退回,1:退回
			l.UpdateReturnByEmail(receiver, err.Error())
		} else {
			l.UpdateReturnByEmail(customer.Email, err.Error())
		}
		l.Logger.Errorf("send mail %v fail: %v", receiver, err)
		return err
	}
	fmt.Println(receiver + " send mail finsh")
	return nil
}

/*
*
更新email状态退回
*/
func (l *AutoMailLogic) UpdateReturnByEmail(emails string, note string) {
	searchContact, err := l.svcCtx.SearchContact.FindOneByEmail(l.ctx, emails)
	if err != nil {
		l.Logger.Errorf(" UpdateReturnByEmail:%v", err)
		return
	}
	if searchContact != nil {
		fmt.Println("系统存在:" + emails + "更新状态为退回")
		//系统退回0:未退回,1:退回
		searchContact.IsReturn = 1
		searchContact.IsSend = 2
		searchContact.Note = note
		err := l.svcCtx.SearchContact.Update(l.ctx, searchContact)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
func (l *AutoMailLogic) ReceiveEmail(emailProviders model.EmailProviders) {

	pop3Server := emailProviders.PopServer
	port := emailProviders.PopPort
	username := emailProviders.Username
	password := emailProviders.Password
	addr := fmt.Sprintf("%s:%s", pop3Server, port)
	// 建立TLS连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{})
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
	for i := count; i > 0; i-- {
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
		//fmt.Println(from)
		if strings.Contains(from, "PostMaster@qq.com") || strings.Contains(from, "Postmaster@163.com") {
			//fmt.Printf("Message %d Text Body: %s\n", i, env.Text)
			// 电子邮件地址正则表达式
			emailRegex := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`

			// 编译正则表达式
			re := regexp.MustCompile(emailRegex)
			emails := re.FindString(env.Text)
			l.UpdateReturnByEmail(emails, "")

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
