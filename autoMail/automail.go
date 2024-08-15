package main

import (
	"automail/autoMail/internal/config"
	"automail/autoMail/internal/logic"
	"automail/autoMail/internal/svc"
	"context"
	"flag"
	"github.com/robfig/cron/v3"

	"github.com/zeromicro/go-zero/core/conf"
)

// var configFile = flag.String("f", "/home/dave/www/autoSendmail/autoMail/etc/config.yaml", "the config file")
var configFile = flag.String("f", "autoMail/etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	cxt := context.Background()
	svcCtx := svc.NewServiceContext(c)

	l := logic.NewAutoMailLogic(cxt, svcCtx)

	crondtask := cron.New(cron.WithSeconds())
	//// 每周二 11:00:00 触发
	_, err := crondtask.AddFunc("0 0 11 * * 2", l.AutoMail)
	//_, err := crondtask.AddFunc("0 21 16 * * 4", l.AutoMail)
	if err != nil {
		return
	}
	crondtask.Start()
	defer crondtask.Stop()
	select {}

}
