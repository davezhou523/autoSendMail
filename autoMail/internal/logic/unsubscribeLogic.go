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

	return
}
