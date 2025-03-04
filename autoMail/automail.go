package main

import (
	"automail/autoMail/internal/config"
	"automail/autoMail/internal/handler"
	"automail/autoMail/internal/logic"
	"automail/autoMail/internal/svc"
	"context"
	"flag"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "autoMail/etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 初始化日志
	logx.MustSetup(c.Log)
	defer logx.Close()
	cxt := context.Background()
	fmt.Println(cxt)
	svcCtx := svc.NewServiceContext(c)
	crondtask := cron.New(cron.WithSeconds())
	////// 每周二 11:00:00 触发
	emailProvidersL := logic.NewEmailProvidersLogic(cxt, svcCtx)
	_, err := crondtask.AddFunc("0 0 0 * * *", emailProvidersL.ResetCountAndTime)
	if err != nil {
		_ = fmt.Errorf("crondtask:%v\n", err)
	}
	l := logic.NewAutoMailLogic(cxt, svcCtx)
	//l.AutoMail()

	_, err = crondtask.AddFunc("0 40 9 * * 2", l.AutoMail)
	if err != nil {
		_ = fmt.Errorf("crondtask:%v\n", err)
	}
	crondtask.Start()
	defer crondtask.Stop()

	//l.CustomizeSend()
	// 启动服务
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	handler.RegisterHandlers(server, svcCtx)
	fmt.Println(svcCtx.Config.Host, svcCtx.Config.Port)
	//fmt.Println(helper.GenerateToken("271416962@qq.com", svcCtx.Config.Secret))
	server.Start()

}
