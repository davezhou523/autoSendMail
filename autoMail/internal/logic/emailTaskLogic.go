package logic

import (
	"automail/autoMail/internal/svc"
	"automail/model"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type EmailTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailTaskLogic {
	return &EmailTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *EmailTaskLogic) worker(wg *sync.WaitGroup, tasks chan *model.SearchContact, providers []*model.EmailProviders) {
	defer wg.Done()

	providerIndex := 0
	for task := range tasks {
		if len(providers) == 0 {
			fmt.Println("没有可用的邮件服务商，任务暂停")
			break
		}

		provider := providers[providerIndex]
		fmt.Sprintf("provider:%v\n", provider.Username)
		fmt.Sprintf("task Email:%v\n", task.Email)

		//err := email.SendMail(provider, task.Recipient, task.Subject, task.Body)
		//if err != nil {
		//	fmt.Printf("[✖] 邮件发送失败: %s\n", err)
		//	updateTaskStatus(db, task.ID, "failed")
		//} else {
		//	updateTaskStatus(db, task.ID, "sent")
		//}

		time.Sleep(1 * time.Second)                          // 防止过快发送
		providerIndex = (providerIndex + 1) % len(providers) // 轮询选择 SMTP 账号
	}
}

func (l *EmailTaskLogic) saveEmailTask(customer *model.SearchContact, emailContent *model.EmailContent) (id int64, err error) {
	emailTask := new(model.EmailTask)
	emailTask.Email = customer.Email
	emailTask.ContentId = emailContent.Id
	emailTask.SendTime = time.Now().Unix()
	et, err := l.svcCtx.EmailTask.Insert(l.ctx, emailTask)
	if err != nil {
		return 0, err
	}
	id, err = et.LastInsertId()

	return id, err
}
