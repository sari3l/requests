package tools

import (
	"fmt"
	"github.com/sari3l/requests"
	"path"
)

func Download(session *requests.Session, rawUrl *string, savePath *string) {
	resp := session.Get(*rawUrl)
	if resp == nil {
		fmt.Printf("[x] 下载失败：%s\n", *rawUrl)
		return
	}
	tempPath := path.Join(*savePath, path.Base(resp.Request.URL.Path))
	err := resp.Save(tempPath)
	if err != nil {
		fmt.Printf("[x] 下载失败：%s\n", *rawUrl)
	} else {
		fmt.Printf("[v] 下载成功：%s，保存路径：%s\n", *rawUrl, tempPath)
	}
}
