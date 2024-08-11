package logic

import (
	"context"

	"automail/autoMail/internal/svc"
	"automail/autoMail/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AutoMailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAutoMailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AutoMailLogic {
	return &AutoMailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AutoMailLogic) AutoMail(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
