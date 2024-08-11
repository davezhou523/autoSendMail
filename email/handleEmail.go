package email

import (
	"automail/db"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Attach_struct struct {
	file_name string
	file_path string
}

func EmailRun() {
	rows, err := db.DB.Query("SELECT id,title,content,attach_id FROM  email_content ")
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		var content string
		var attach_id string
		if err := rows.Scan(&id, &title, &content, &attach_id); err != nil {
			log.Fatalf("Scan 失败: %v", err)
		}

		attch, err := getAttch(attach_id)
		if err != nil {
			return
		}
		go ScheduleEmail(1*time.Second, content, title, &attch)
		//fmt.Printf("email_content: %d, %s\n", id, attach_id)
	}

	//go email.ScheduleEmail(5*24*time.Hour, "content2.txt", "Content 2")
	//var recipients = []string{"a@gmail.com", "b@gmail.com"}
	//for key, receiver := range recipients {
	//	fmt.Println(key, receiver)
	//}
	// 保持程序运行
	//select {}
	//time.Sleep(10 * time.Second)
}
func getAttch(attach_id string) ([]Attach_struct, error) {
	var ids []int
	err := json.Unmarshal([]byte(attach_id), &ids)
	if err != nil {
		fmt.Printf("attach_id json失败: %v\n", err)
		return nil, err
	}
	// 构建查询语句
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?" // 占位符
		args[i] = id          // 参数
	}
	query := fmt.Sprintf("SELECT  file_name,file_path  FROM  attach where id IN (%s)", strings.Join(placeholders, ", "))
	rows, err := db.DB.Query(query, args...)
	defer rows.Close()
	if err != nil {
		fmt.Printf("attach查询失败: %v\n", err)
		return nil, err
	}
	attachArr := make([]Attach_struct, 0)
	for rows.Next() {
		attach := Attach_struct{}
		if err := rows.Scan(&attach.file_name, &attach.file_path); err != nil {
			log.Fatalf("Scan 失败: %v", err)
		}
		attachArr = append(attachArr, attach)
		fmt.Printf("attach: %s, %s\n", attach.file_name, attach.file_path)
	}
	return attachArr, nil

}
