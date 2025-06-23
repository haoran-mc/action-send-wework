package model

import (
	"fmt"

	"github.com/haoran-mc/golib/pkg/server/http"
)

// refer https://work.weixin.qq.com/api/doc/90000/90136/91770
type Bot struct {
	Key string
}

type BotRequest struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text,omitempty"`
}

type BotResponse struct {
	ErrorCode    int64  `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

func (r *Bot) CreateBaseURL() string {
	return fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", r.Key)
}

func (r *Bot) SendText(text string) (res *BotResponse, err error) {
	data := BotRequest{
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
		return &BotResponse{
			ErrorCode:    -1,
			ErrorMessage: "fail to request",
		}, err
	}
	return
}
