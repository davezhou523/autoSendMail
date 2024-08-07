package email

import "time"

func EmailRun() {
	var filename1 string = "static/content1.html"
	go ScheduleEmail(1*time.Second, filename1, "Content 1")
	//go email.ScheduleEmail(5*24*time.Hour, "content2.txt", "Content 2")
	//var recipients = []string{"a@gmail.com", "b@gmail.com"}
	//for key, receiver := range recipients {
	//	fmt.Println(key, receiver)
	//}
	// 保持程序运行
	//select {}
	time.Sleep(10 * time.Second)
}
