package main

import (
	"automail/mail"
	"time"
)

func main() {
	go mail.ScheduleEmail(3*24*time.Hour, "content1.txt", "Content 1")
	go mail.ScheduleEmail(5*24*time.Hour, "content2.txt", "Content 2")
	// 保持程序运行
	select {}
}
