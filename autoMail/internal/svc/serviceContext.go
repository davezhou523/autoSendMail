package svc

import (
	"automail/autoMail/internal/config"
	"database/sql"
	"fmt"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	db, err := sql.Open("mysql", c.DataSource.DataSourceName)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}
	db.SetMaxOpenConns(c.DataSource.MaxOpenConns)
	db.SetMaxIdleConns(c.DataSource.MaxIdleConns)

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
