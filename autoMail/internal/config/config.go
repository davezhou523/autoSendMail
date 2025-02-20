package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	//Log        logx.LogConf
	DataSource struct {
		DataSourceName string
	}
	Unsubscribe string
	ReplyTo     string
	SmtpSource  struct {
		Server   string
		Port     int
		Username string
		Password string
	}
	PopSource struct {
		Server   string
		Port     int
		Username string
		Password string
	}
}
