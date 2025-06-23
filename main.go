package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	googledrive "github.com/haoran-mc/action-send-wework/internal/repository/google-drive"
	"github.com/haoran-mc/action-send-wework/internal/service"
	"github.com/haoran-mc/golib/pkg/env"
	"github.com/haoran-mc/golib/pkg/log"
	"google.golang.org/api/drive/v3"
)

var readfiles = []string{"ideas.txt", "daily-reminder.txt", "memorial-days.txt"}

var fileRandomProbabilityMap = map[string]float32{
	"ideas.txt":          1, // plain text
	"daily-reminder.txt": 1, // plain text
	"memorial-days.txt":  1, // formatting text
}

func main() {
	var err error
	for range 5 {
		time.Sleep(10 * time.Second)

		credentialsJSON := env.GetEnv("GDRIVE_CREDENTIALS", "")
		err = googledrive.InitDriveService(credentialsJSON)
		if err != nil {
			log.Error("fail to init drive service", err)
			continue
		}

		// 1. read googledrive shared directory
		dirID := env.GetEnv("DIR_ID", "")
		files, err := googledrive.ReadDir(dirID)
		if err != nil {
			log.Error("fail to read shared directory", err)
			continue
		}

		// 2. read files and generate the send str
		sendStr := generateSendStr(files)
		if len(sendStr) == 0 {
			sendStr = "青山落日，秋月春风。当真是朝如青丝暮成雪，是非成败转头空。"
		}

		// 3. wework bot send
		botKey := env.GetEnv("BOT_KEY", "")
		err = service.BotSend(botKey, sendStr)
		if err != nil {
			log.Error("bot send message failed", err)
			continue
		}
		break
	}
	log.Info("successfully")
}

func generateSendStr(files []*drive.File) string {
	msgs := []string{}
	for _, filename := range readfiles { // 按顺序
		for _, f := range files {
			if f.Name != filename {
				continue
			}
			// 1. 是否读此文件？
			rf := rand.Float32()
			// log.Info(fmt.Sprintf("%s\t%s\t%s: %b\n", f.Id, f.Name, f.Md5Checksum, rf))
			if rf > fileRandomProbabilityMap[filename] {
				continue
			}
			fileContent, err := googledrive.ReadFile(f.Id)
			if err != nil {
				log.Error("fail to read file", err)
				msgs = append(msgs, fmt.Sprintf("read f.Name failed, error: %v\n", err))
				continue
			}
			// 2. 不同文件，不同格式，读数据方式不同
			switch filename {
			case "ideas.txt":
				msgs = append(msgs, service.RandomLine(fileContent))
			case "daily-reminder.txt":
				msgs = append(msgs, service.RandomLine(fileContent))
			case "memorial-days.txt":
				msgs = append(msgs, service.ReadFormattingText(fileContent)...)
			}
		}
	}
	return strings.Join(msgs, "\n")
}
