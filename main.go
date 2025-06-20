package main

import (
	"fmt"
	"math/rand"

	"github.com/haoran-mc/action-send-wework/internal/model"
	googledrive "github.com/haoran-mc/action-send-wework/internal/repository/google-drive"
	"github.com/haoran-mc/golib/pkg/env"
	"github.com/haoran-mc/golib/pkg/log"
	"google.golang.org/api/drive/v3"
)

var readfiles = []string{"ideas.txt", "memorial-days.txt", "daily-reminder.txt"}

var fileRandomProbabilityMap = map[string]float32{
	"ideas.txt":          1, // plain text
	"memorial-days.txt":  0, // formatting text
	"daily-reminder.txt": 1, // plain text
}

func main() {
	var err error
	err = googledrive.InitDriveService()
	if err != nil {
		// TODO 重试
	}

	// 1. read googledrive shared directory
	dirID := env.GetEnv("DIR_ID", "")
	files, err := googledrive.ReadDir(dirID)
	if err != nil {
		// TODO 重试
	}

	// 2. read files and generate the send str
	sendStr := generateSendStr(files)
	if len(sendStr) == 0 {
		sendStr = "青山落日，秋月春风。当真是朝如青丝暮成雪，是非成败转头空。"
	}

	// 3. wework bot send
	err = send(sendStr)
	if err != nil {
		// TODO 重试
	}
}

func generateSendStr(files []*drive.File) (sendStr string) {
	for _, filename := range readfiles {
		for _, f := range files {
			if f.Name != filename {
				continue
			}
			rf := rand.Float32()
			log.Info(fmt.Sprintf("%s\t%s\t%s: %b\n", f.Id, f.Name, f.Md5Checksum, rf))
			if rf > fileRandomProbabilityMap[filename] {
				continue
			}
			fileContent, err := googledrive.ReadFile(f.Id)
			if err != nil {
				log.Error("fail to read file", err)
				sendStr += fmt.Sprintf("read f.Name failed, error: %v\n", err)
				continue
			}
			sendStr += string(fileContent)
		}
	}
	return
}

func send(text string) error {
	robotKey := env.GetEnv("BOT_KEY", "")
	robot := model.Robot{Key: robotKey}

	res, err := robot.SendText(text)
	if err != nil {
		log.Error("fail to send text", err)
	} else if res != nil && res.ErrorCode != 0 {
		log.Error("fail to send text", res.ErrorMessage)
		err = fmt.Errorf("fail to send text, " + res.ErrorMessage)
	}
	return err
}
