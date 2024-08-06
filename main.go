package main

import (
	"automail/email"
	"fmt"
	"log"
	"os"
	"time"
)

//var Loger *log.Logger

func init() {
	file := "./log/" + time.Now().Format("2006-01-02") + ".txt"
	//os.MkdirAll(file, 0766)
	fmt.Println(file)
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {

		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[automail]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	//Loger = log.New(logFile, "[automail]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
	return
}
func main() {
	log.Println("start send mail")
	var filename1 string = "static/content1.html"
	go email.ScheduleEmail(1*time.Second, filename1, "Content 1")
	//go email.ScheduleEmail(5*24*time.Hour, "content2.txt", "Content 2")
	//var recipients = []string{"a@gmail.com", "b@gmail.com"}
	//for key, receiver := range recipients {
	//	fmt.Println(key, receiver)
	//}
	// 保持程序运行
	//select {}
	time.Sleep(3 * time.Second)
}
