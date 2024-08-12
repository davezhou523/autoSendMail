package svc

import (
	"automail/autoMail/internal/config"
	"automail/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	Attach        model.AttachModel
	EmailContent  model.EmailContentModel
	SearchContact model.SearchContactModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.DataSource.DataSourceName)
	return &ServiceContext{
		Config:        c,
		Attach:        model.NewAttachModel(conn),
		EmailContent:  model.NewEmailContentModel(conn),
		SearchContact: model.NewSearchContactModel(conn),
	}
}
