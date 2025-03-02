package logic

import (
	"automail/autoMail/internal/svc"
	"automail/model"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
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
