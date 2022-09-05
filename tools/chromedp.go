package tools

import (
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/sari3l/requests"
	"io/ioutil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

// 将不常用的但又有需求的执行内容放在这里，避免Response过于庞杂

// FileDownload 文件下载器，参考 https://github.com/chromedp/examples/blob/master/download_image/main.go
func FileDownload(response *requests.Response, fileTypes []string, savePath string) string {
	var typeDict = make(map[string]bool, 0)
	if fileTypes == nil || len(fileTypes) == 0 {
		return ""
	} else {
		for _, fileType := range fileTypes {
			typeDict[strings.ToLower(fileType)] = true
		}
	}
	if len(savePath) == 0 {
		savePath, _ = ioutil.TempDir("", strconv.FormatInt(time.Now().Unix(), 10))
		fmt.Printf("[*] 未指定保存路径，创建临时文件夹：%s\n", savePath)
	}
	listener := downloadNetworkEventListener(response.Session, &typeDict, &savePath)
	response.CustomRender(listener, nil, nil)
	return savePath
}

func downloadNetworkEventListener(session *requests.Session, typeDict *map[string]bool, savePath *string) func(ev interface{}) {
	return func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			fileURL := ev.Request.URL
			if fileURLParser, err := url.Parse(fileURL); err == nil {
				pos := strings.LastIndex(fileURLParser.Path, ".")
				if pos != -1 && (*typeDict)[fileURLParser.Path[pos+1:len(fileURLParser.Path)]] == true {
					//fmt.Println(fileURL)
					go download(session, &fileURL, savePath)
				}
			}
		}
	}
}

func download(session *requests.Session, rawUrl *string, savePath *string) {
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
