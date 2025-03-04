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
	Unsubscribe    string
	ReplyTo        string
	UnsubscribeAPI string
	Secret         string
	EmailContentId uint64
}
