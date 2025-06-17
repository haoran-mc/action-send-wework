package googledrive

import (
	"context"
	"io"

	"github.com/haoran-mc/golib/pkg/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func ReadFile(credentialsJSON, fileID string) (fileContent []byte, err error) {
	ctx := context.Background()

	// 从环境变量获取凭证和配置
	if credentialsJSON == "" {
		log.Error("empty environment variable GDRIVE_CREDENTIALS")
		return nil, err
	}

	// Google Drive 中文件的 ID
	if fileID == "" {
		log.Error("empty environment variable FILE_ID")
		return nil, err
	}

	// 使用服务账号凭证创建授权的客户端
	creds, err := google.CredentialsFromJSON(ctx, []byte(credentialsJSON), drive.DriveReadonlyScope)
	if err != nil {
		log.Error("unable to create client from credentials", err)
		return nil, err
	}

	// 创建一个经过授权的 Drive 服务
	driveService, err := drive.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Error("unable to create Drive service", err)
		return nil, err
	}

	// 下载文件内容
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
	return
}
