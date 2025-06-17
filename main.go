package main

import (
	"github.com/haoran-mc/action-send-wework/model"
	googledrive "github.com/haoran-mc/action-send-wework/repository/google-drive"
	"github.com/haoran-mc/golib/pkg/env"
	"github.com/haoran-mc/golib/pkg/log"
)

func main() {
	robotKey := env.GetEnv("BOT_KEY", "")
	robot := model.Robot{Key: robotKey}

	credentialsJSON := env.GetEnv("GDRIVE_CREDENTIALS", ``)
	fileID := env.GetEnv("FILE_ID", "")

	fileContent, err := googledrive.ReadFile(credentialsJSON, fileID)
	if err != nil {
		// 重试
	}

	// 文本
	res, err := robot.SendText(string(fileContent))
	if err != nil {
		log.Error("fail to send text", err)
	} else if res != nil && res.ErrorCode != 0 {
		log.Error("fail to send text", res.ErrorMessage)
	}
}
