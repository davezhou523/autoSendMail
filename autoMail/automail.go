package main

import (
	"automail/autoMail/internal/config"
	"automail/autoMail/internal/logic"
	"automail/autoMail/internal/svc"
	"context"
	"flag"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"regexp"
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
	//l.ConvertEmailDomainLower()
	l.CustomizeSend()
	//l.UpdateReturnByEmail("pgfilters@premiumguard.com")
	//l.ReceiveEmail()
	//vali()
	crondtask := cron.New(cron.WithSeconds())
	////// 每周二 11:00:00 触发
	_, err := crondtask.AddFunc("0 08 11 * * 2", l.AutoMail)
	//_, err := crondtask.AddFunc("*/10 * * * * *", l.AutoMail)

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

func vali() {
	text := "无法发送到 ann@arthron5.com 退信原因 收件人邮件地址（ann@arthron5.com）不存在，邮件无法送达。\nhost arthron5.com[66.102.132.129] said: 550 No Such User Here (in reply to RCPT TO command) 解决方案 请联系您的收件人，重新核实邮箱地址，或发送到其他收信邮箱。 您也可以 向管理员报告此\n信 ( http://mail.qq.com/cgi-bin/help_feedback_person?sender=noratf@foxmail.com&receiver=ann@arthron5.com&sendtime=2024-08-29 11:20:29&reason=host arthron5.com[66.102.132.129] said:   550 No Such User Here (in reply to RCPT TO command)&subject=Certified+Quality%3A+Trust+Our+ISO+and+CE+Certified+Gloves ) 。\n\n此外，您还可以 点击这里 ( http://service.mail.qq.com/cgi-bin/help?subtype=1&&id=29&&no=188 ) 获取更多关于退信的帮助信息。"
	//emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`
	// 编译正则表达式
	re := regexp.MustCompile(emailRegex)

	emails := re.FindString(text)
	fmt.Println("emails:" + emails)

}
