package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DataSource struct {
		DataSourceName string
		MaxOpenConns   int
		MaxIdleConns   int
	}
}
