package svc

import (
	"automail/autoMail/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	DB     sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	//db, err := sql.Open("mysql", c.DataSource.DataSourceName)
	db := sqlx.NewMysql(c.DataSource.DataSourceName)
	//if err != nil {
	//	fmt.Println(err)
	//	//panic(fmt.Sprintf("Failed to connect to database: %s", err))
	//}
	//db.SetMaxOpenConns(c.DataSource.MaxOpenConns)
	//db.SetMaxIdleConns(c.DataSource.MaxIdleConns)

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
