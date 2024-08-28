package main

import (
	"automail/autoMail/internal/config"
	"automail/autoMail/internal/logic"
	"automail/autoMail/internal/svc"
	"context"
	"flag"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// var configFile = flag.String("f", "/home/dave/www/autoSendmail/autoMail/etc/config.yaml", "the config file")
var configFile = flag.String("f", "autoMail/etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 初始化日志
	logx.MustSetup(c.Log)
	defer logx.Close()
	cxt := context.Background()
	svcCtx := svc.NewServiceContext(c)
	//
	l := logic.NewAutoMailLogic(cxt, svcCtx)
	//l.AutoMail()
	crondtask := cron.New(cron.WithSeconds())
	////// 每周二 11:00:00 触发
	_, err := crondtask.AddFunc("0 31 10 * * 4", l.AutoMail)
	//_, err := crondtask.AddFunc("*/10 * * * * *", l.AutoMail)

	//crondtask.AddFunc("0 48 13 * * 4", l.AutoMail)
	//crondtask.AddFunc("0 52 13 * * 4", l.AutoMail)
	if err != nil {
		l.Logger.Errorf("crondtask:%v\n", err)
		return
	}
	crondtask.Start()
	defer crondtask.Stop()
	//select {}

	// 启动服务
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Start()

}
