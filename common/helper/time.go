package helper

import (
	"math/rand"
	"time"
)

func StrtimeToInt(datetime string, timeLayout string) int64 {
	//日期转化为时间戳
	if timeLayout == "" {
		timeLayout = "2006-01-02 15:04:05" //转化所需模板
	}
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix() //转化为时间戳 类型是int64
	return timestamp
}

// 生成随机字符串
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
