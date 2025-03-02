package logic

import (
	"automail/autoMail/internal/svc"
	"automail/model"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type EmailProvidersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailProvidersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailProvidersLogic {
	return &EmailProvidersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailProvidersLogic) getProvidersList(user_id int64, company_id int64) (emailProviders []*model.EmailProviders, err error) {
	// 查询可用服务商
	providers, _ := l.svcCtx.EmailProviders.FindAll(l.ctx, user_id, company_id)
	if len(providers) == 0 {
		msg := "未配置发送邮件"
		l.Logger.Infof(msg)
		return emailProviders, fmt.Errorf(msg)
	}
	for _, p := range providers {
		if p.DailyLimit < p.SentCount {
			continue
		}
		emailProviders = append(emailProviders, p)
	}
	if len(emailProviders) == 0 {
		msg := "超出邮件每日限额"
		l.Logger.Infof(msg)
		return emailProviders, fmt.Errorf(msg)
	}
	return emailProviders, err
}
func (l *EmailProvidersLogic) updateProvider(user_id int64, company_id int64) (*model.EmailProviders, error) {

	// 查询可用服务商
	_, _ = l.svcCtx.EmailProviders.FindAll(l.ctx, user_id, company_id)

	return nil, fmt.Errorf("no available providers")
}
