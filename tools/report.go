package tools

import (
	"encoding/base64"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/sari3l/requests"
	"html/template"
	"io/ioutil"
	"time"
)

// 参考 https://github.com/sensepost/gowitness

func ReportHTTPResponse(response *requests.Response, savePath *string) string {
	t, err := template.ParseFiles("template/report.tpl")
	if err != nil {
		fmt.Println(err)
	}
	if savePath == nil {
		savePath = new(string)
	}
	tmp, err := ioutil.TempFile(*savePath, "*.html")
	if err != nil {
		fmt.Println(err)
	}
	consoleLogs := new([]ConsoleLog)
	networkLogs := map[string]NetworkLog{}
	buf := new([]byte)
	// 增加监听事件
	_ = response.CustomRender([]func(ev interface{}){EventListenerConsoleAPICalled(consoleLogs), EventListenerNetwork(&networkLogs)}, nil, chromedp.FullScreenshot(buf, 100))
	data := map[string]any{
		"URL":         response.Request.URL.String(),
		"StatusCode":  response.StatusCode,
		"Snapshot":    base64.StdEncoding.EncodeToString(*buf),
		"Datetime":    time.Now().String(),
		"ConsoleLogs": consoleLogs,
		"NetworkLogs": networkLogs,
		"Title":       response.Title(),
		"Header":      response.Header,
		"DOM":         response.Html,
	}
	if response.Request.URL.Scheme == "https" {
		data["Certificates"] = response.TLS.PeerCertificates
	}
	t.Execute(tmp, data)
	return tmp.Name()
}
