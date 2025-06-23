package googledrive

import (
	"context"
	"fmt"
	"io"

	"github.com/haoran-mc/golib/pkg/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var driveService *drive.Service

func InitDriveService(credentialsJSON string) (err error) {
	if driveService != nil {
		return nil
	}
	if credentialsJSON == "" {
		log.Error("empty environment variable GDRIVE_CREDENTIALS")
		return fmt.Errorf("credentials JSON is empty")
	}

	ctx := context.Background()

	// 使用服务账号凭证创建授权的客户端
	creds, err := google.CredentialsFromJSON(ctx, []byte(credentialsJSON), drive.DriveReadonlyScope)
	if err != nil {
		log.Error("unable to create client from credentials" + err.Error())
		return
	}
	// 创建经过授权的 Drive 服务
	driveService, err = drive.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Error("fail to create drive service", err)
	}
	return
}

var fileMap = map[string][]byte{}

func ReadFile(fileID string) (fileContent []byte, err error) {
	if fileID == "" {
		log.Error("FILE_ID not provided")
		return nil, err
	}

	if fileContent, ok := fileMap[fileID]; ok {
		return fileContent, nil
	}

	resp, err := driveService.Files.Get(fileID).Download()
	if err != nil {
		log.Error("unable to download file", err)
		return nil, err
	}
	defer resp.Body.Close()

	fileContent, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error("fail to read file content", err)
	}
	fileMap[fileID] = fileContent
	return
}

var driveFiles []*drive.File

func ReadDir(dirID string) ([]*drive.File, error) {
	if driveFiles != nil {
		return driveFiles, nil
	}
	if dirID == "" {
		log.Error("empty environment variable DIR_ID")
		return nil, fmt.Errorf("directory ID is empty")
	}

	// 'trashed = false' 不列出回收站中文件
	query := fmt.Sprintf("'%s' in parents and trashed = false", dirID)
	pageToken := ""

	// drive api 使用分页来返回结果，循环获取所有页面
	for {
		req := driveService.Files.List().
			Q(query).
			Fields("nextPageToken, files(*)"). // drive api 默认返回 id, name, mimeType，使用 "*" 请求所有可用字段
			PageToken(pageToken)

		r, err := req.Do()
		if err != nil {
			log.Error("unable to retrieve files", err)
			return nil, err
		}

		driveFiles = append(driveFiles, r.Files...)

		// 如果有下一页，则更新 pageToken 继续
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return driveFiles, nil
}
