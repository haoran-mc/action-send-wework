package model

import (
	"fmt"

	"github.com/haoran-mc/golib/pkg/server/http"
)

// refer https://work.weixin.qq.com/api/doc/90000/90136/91770
type Robot struct {
	Key string
}

// 机器人接口请求
type RobotRequest struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text,omitempty"`
}

// 机器人接口响应
type RobotResponse struct {
	ErrorCode    int64  `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

// 拼接地址
func (r *Robot) CreateBaseURL() string {
	return fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", r.Key)
}

// 发送纯文本
func (r *Robot) SendText(text string) (res *RobotResponse, err error) {
	data := RobotRequest{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: text,
		},
	}
	baseURL := r.CreateBaseURL()
	err = http.PostJson(baseURL, data)
	if err != nil {
		return &RobotResponse{
			ErrorCode:    -1,
			ErrorMessage: "fail to request",
		}, err
	}
	return
}
