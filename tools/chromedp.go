package tools

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/sari3l/requests"
	"github.com/sari3l/requests/types"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ConsoleLog struct {
	Type  string
	Value string
}

type NetworkLog struct {
	RequestID  string
	Time       time.Time
	URL        string
	Type       string
	StatusCode int64
	EndURL     string
	IP         string
	Error      string
}

// 将不常用的但又有需求的执行内容放在这里，避免Response过于庞杂

// FileDownload 文件下载器，参考 https://github.com/chromedp/examples/blob/master/download_image/main.go
func FileDownload(response *requests.Response, fileTypes []string, savePath string) string {
	return CustomDownloader(response, fileTypes, savePath, nil, nil)
}

func CustomDownloader(response *requests.Response, fileTypes []string, savePath string, flags []chromedp.ExecAllocatorOption, actions ...chromedp.Action) string {
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
	listener := EventListenerParserToDownload(response.Session, &typeDict, &savePath)
	response.CustomRender([]func(ev interface{}){listener}, flags, actions...)
	return savePath
}

func EventListenerParserToDownload(session *requests.Session, typeDict *map[string]bool, savePath *string) func(ev interface{}) {
	return func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			fileURL := ev.Request.URL
			if fileURLParser, err := url.Parse(fileURL); err == nil {
				pos := strings.LastIndex(fileURLParser.Path, ".")
				if pos != -1 && (*typeDict)[fileURLParser.Path[pos+1:len(fileURLParser.Path)]] == true {
					//fmt.Println(fileURL)
					go Download(session, &fileURL, savePath)
				}
			}
		}
	}
}

func EventListenerNetwork(networkLogs *map[string]NetworkLog) func(ev interface{}) {
	return func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			(*networkLogs)[string(ev.RequestID)] = NetworkLog{
				RequestID: string(ev.RequestID),
				Type:      "HTTP",
				Time:      time.Time(*ev.Timestamp),
				URL:       ev.Request.URL,
			}
		case *network.EventResponseReceived:
			if log, ok := (*networkLogs)[string(ev.RequestID)]; ok {
				log.StatusCode = ev.Response.Status
				log.EndURL = ev.Response.URL
				log.IP = ev.Response.RemoteIPAddress
				(*networkLogs)[string(ev.RequestID)] = log
			}
		case *network.EventLoadingFailed:
			if log, ok := (*networkLogs)[string(ev.RequestID)]; ok {
				log.Error = ev.ErrorText
				(*networkLogs)[string(ev.RequestID)] = log
			}
		case *network.EventWebSocketCreated:
			(*networkLogs)[string(ev.RequestID)] = NetworkLog{
				RequestID: string(ev.RequestID),
				Type:      "Socket",
				URL:       ev.URL,
			}

		case *network.EventWebSocketHandshakeResponseReceived:
			if log, ok := (*networkLogs)[string(ev.RequestID)]; ok {
				log.StatusCode = ev.Response.Status
				log.Time = time.Time(*ev.Timestamp)
				(*networkLogs)[string(ev.RequestID)] = log
			}
		case *network.EventWebSocketFrameError:
			if log, ok := (*networkLogs)[string(ev.RequestID)]; ok {
				log.Error = ev.ErrorMessage
				(*networkLogs)[string(ev.RequestID)] = log
			}
		}
	}
}

func EventListenerConsoleAPICalled(result *[]ConsoleLog) func(ev interface{}) {
	return func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			buf := ""
			for _, arg := range ev.Args {
				buf += string(arg.Value)
			}
			*result = append(*result, ConsoleLog{string(ev.Type), buf})
		case *runtime.EventExceptionThrown:
			*result = append(*result, ConsoleLog{"exception", ev.ExceptionDetails.Error()})
		}
	}
}

func ActionFuncSetCookies(domain *string, cookies *types.Dict) chromedp.ActionFunc {
	//cookiesResp := &resp.cookies
	//if useExtCookies || cookiesResp == nil || len(*cookiesResp) == 0 {
	return func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		for name, value := range *cookies {
			err := network.SetCookie(name, value).
				WithExpires(&expr).
				WithDomain(*domain).
				WithPath("/").
				WithHTTPOnly(false).
				WithSecure(false).
				Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
	//} else {
	//return func(ctx context.Context) error {
	//	for _, cookie := range *cookiesResp {
	//		err := network.SetCookie(cookie.Name, cookie.String()).
	//			WithExpires((*cdp.TimeSinceEpoch)(&cookie.Expires)).
	//			WithDomain(cookie.Domain).
	//			WithPath(cookie.Path).
	//			WithHTTPOnly(cookie.HttpOnly).
	//			WithSecure(cookie.Secure).
	//			Do(ctx)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//	return nil
	//}
	//}
}
