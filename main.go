package main

import (
	"fmt"
	"math/rand"

	"github.com/haoran-mc/action-send-wework/model"
	googledrive "github.com/haoran-mc/action-send-wework/repository/google-drive"
	"github.com/haoran-mc/golib/pkg/env"
	"github.com/haoran-mc/golib/pkg/log"
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

	dirID := env.GetEnv("DIR_ID", "")
	files, err := googledrive.ReadDir(dirID)
	if err != nil {
		// TODO 重试
	}

	var sendStr string
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
	log.Info(sendStr)

	robotKey := env.GetEnv("BOT_KEY", "")
	robot := model.Robot{Key: robotKey}

	res, err := robot.SendText(sendStr)
	if err != nil {
		log.Error("fail to send text", err)
	} else if res != nil && res.ErrorCode != 0 {
		log.Error("fail to send text", res.ErrorMessage)
	}
}
