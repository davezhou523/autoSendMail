package logic

import (
	"automail/autoMail/internal/svc"
	"automail/model"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

func (l *EmailTaskLogic) saveData(customer *model.SearchContact, emailContent *model.EmailContent, provider *model.EmailProviders) *model.EmailTask {
	emailTask := new(model.EmailTask)
	emailTask.Email = customer.Email
	emailTask.ContentId = emailContent.Id
	emailTask.SendTime = time.Now().Unix()
	emailTask.UserId = provider.UserId
	emailTask.CompanyId = provider.CompanyId
	emailTask.ProviderEmail = provider.Username
	return emailTask

}
func (l *EmailTaskLogic) saveEmailTask(customer *model.SearchContact, emailContent *model.EmailContent, provider *model.EmailProviders) (id int64, err error) {

	et, err := l.svcCtx.EmailTask.Insert(l.ctx, l.saveData(customer, emailContent, provider))
	if err != nil {
		return 0, err
	}
	id, err = et.LastInsertId()

	return id, err
}
func (l *EmailTaskLogic) saveEmailTaskWithSession(session sqlx.Session, customer *model.SearchContact, emailContent *model.EmailContent, provider *model.EmailProviders) (id int64, err error) {

	EmailTaskSession := l.svcCtx.EmailTask.WithSession(session)
	et, err := EmailTaskSession.Insert(l.ctx, l.saveData(customer, emailContent, provider))
	if err != nil {
		return 0, err
	}
	id, err = et.LastInsertId()

	return id, err
}
