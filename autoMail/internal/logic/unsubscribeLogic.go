package logic

import (
	"automail/autoMail/internal/svc"
	"automail/autoMail/internal/types"
	"automail/common/helper"
	"context"

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
	token := helper.GenerateToken(req.Email, l.svcCtx.Config.Secret)
	if token != req.Token {
		l.Logger.Errorf("token验证失败,email:%v,请求:%v,系统token:%v\n", req.Email, req.Token, token)
		//return nil, fmt.Errorf("token验证失败")
	}
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
