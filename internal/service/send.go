package service

import (
	"fmt"

	"github.com/haoran-mc/action-send-wework/internal/model"
	"github.com/haoran-mc/golib/pkg/log"
)

func BotSend(botKey, text string) error {
	bot := model.Bot{Key: botKey}

	res, err := bot.SendText(text)
	if err != nil {
		log.Error("fail to send text", err)
	} else if res != nil && res.ErrorCode != 0 {
		log.Error("fail to send text", res.ErrorMessage)
		err = fmt.Errorf("fail to send text, " + res.ErrorMessage)
	}
	return err
}
