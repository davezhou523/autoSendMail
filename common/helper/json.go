package helper

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReturnContentStruct struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

func (r *ReturnContentStruct) JsonArrPush(jsonData string) string {
	if len(jsonData) > 0 {
		var rc []ReturnContentStruct
		var rcList []ReturnContentStruct
		err := json.Unmarshal([]byte(jsonData), &rc)
		if err != nil {
			logx.Errorf("jsonArrPush:%v,jsonData：%v", err, jsonData)
		}
		for _, value := range rc {
			rcList = append(rcList, value)
		}
		rcList = append(rcList, *r)
		contentByte, _ := json.Marshal(rcList)
		return string(contentByte)

	} else {
		var rcList []ReturnContentStruct
		rcList = append(rcList, *r)
		returnContentByte, _ := json.Marshal(rcList)
		return string(returnContentByte)
	}
}
