package svc

import (
	"automail/autoMail/internal/config"
	"automail/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	SqlConn        sqlx.SqlConn
	Attach         model.AttachModel
	EmailContent   model.EmailContentModel
	SearchContact  model.SearchContactModel
	EmailTask      model.EmailTaskModel
	EmailProviders model.EmailProvidersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.DataSource.DataSourceName)

	return &ServiceContext{
		Config:         c,
		SqlConn:        conn,
		Attach:         model.NewAttachModel(conn),
		EmailContent:   model.NewEmailContentModel(conn),
		SearchContact:  model.NewSearchContactModel(conn),
		EmailTask:      model.NewEmailTaskModel(conn),
		EmailProviders: model.NewEmailProvidersModel(conn),
	}
}
