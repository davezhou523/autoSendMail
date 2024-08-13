package main

import (
	"automail/autoMail/internal/config"
	"automail/autoMail/internal/logic"
	"automail/autoMail/internal/svc"
	"context"
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "/home/dave/www/autoSendmail/autoMail/etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	cxt := context.Background()
	svcCtx := svc.NewServiceContext(c)
	l := logic.NewAutoMailLogic(cxt, svcCtx)
	l.AutoMail()
	select {}
	//server := rest.MustNewServer(c.RestConf)
	//defer server.Stop()

	//handler.RegisterHandlers(server, ctx)

	//fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	//server.Start()
}
