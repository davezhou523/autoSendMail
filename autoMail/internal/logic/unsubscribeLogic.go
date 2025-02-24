package logic

import (
	"context"

	"automail/autoMail/internal/svc"
	"automail/autoMail/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsubscribeLogic {
	return &UnsubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsubscribeLogic) Unsubscribe(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	searchContact, err := l.svcCtx.SearchContact.FindOneByEmail(l.ctx, req.Email)
	if err != nil {
		return nil, err
	}
	searchContact.IsSend = 2
	searchContact.Note = "用户取消订阅"
	err = l.svcCtx.SearchContact.Update(l.ctx, searchContact)
	if err != nil {
		return nil, err
	}

	return &types.Response{Code: 0, Msg: "success"}, err
}
