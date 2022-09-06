package requests

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"github.com/sari3l/requests/parser"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"
)

type Response struct {
	*http.Response
	Session *Session
	cookies []*http.Cookie
	History []*Response
	Html    string
	Ok      bool
	Raw     []byte
	Time    int64
}

func (resp *Response) Json() *gjson.Result {
	if resp == nil {
		return nil
	}
	g := gjson.Parse(resp.Html)
	return &g
}

func (resp *Response) XPath() *parser.XpathNode {
	return parser.XpathParser(&resp.Html)
}

func (resp *Response) Text() string {
	text := ""

	domDoc := html.NewTokenizer(strings.NewReader(resp.Html))
	previousStartToken := domDoc.Token()
loopDom:
	for {
		tt := domDoc.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDom // End of the document,  done
		case tt == html.StartTagToken:
			previousStartToken = domDoc.Token()
		case tt == html.TextToken:
			if previousStartToken.Data == "script" {
				continue
			}
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDoc.Text())))
			if len(TxtContent) > 0 {
				text += TxtContent + "\n"
			}
		}
	}
	return text
}

func (resp *Response) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(resp.Raw)
	f.Sync()
	return err
}

func (resp *Response) ContentType() string {
	return resp.Response.Header.Get("Content-Type")
}

func (resp *Response) Title() string {
	nodes := resp.XPath().Find("/html/head/title")
	if len(nodes) == 1 {
		return nodes[0].Text()
	}
	return ""
}

func (resp *Response) URLs() []string {
	links := linkRegexCompiled.FindAllString(resp.Html, -1)
	originUrl := resp.Request.URL
	return *processLinks(originUrl, &links)
}

// render chromeless 页面渲染
func (resp *Response) render(targetListenerFunctions []func(ev interface{}), customFlags []chromedp.ExecAllocatorOption, tasks ...chromedp.Action) *Response {
	var flags = []chromedp.ExecAllocatorOption{
		chromedp.Flag("ignore-certificate-errors", !resp.Session.Verify),
		chromedp.Flag("headless", true),
		chromedp.Flag("enable-automation", false),
		chromedp.ProxyServer(resp.Session.Proxy),
		chromedp.UserAgent(resp.Request.UserAgent()),
	}
	flags = append(flags, customFlags...)

	opts := append(chromedp.DefaultExecAllocatorOptions[:], flags...)
	chromeCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	tabCtx, cancel := chromedp.NewContext(chromeCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	// 检查是否启动
	if err := chromedp.Run(tabCtx); err != nil {
		panic(err)
	}

	if targetListenerFunctions != nil && len(targetListenerFunctions) > 0 {
		for _, targetListenerFunc := range targetListenerFunctions {
			chromedp.ListenTarget(tabCtx, targetListenerFunc)
		}
	}

	actions := []chromedp.Action{
		actionBypassWebDriver(),
		actionFuncSetCookies(resp),
		chromedp.Navigate(resp.Request.URL.String()),
	}

	for _, task := range tasks {
		if task != nil {
			actions = append(actions, task)
		}
	}
	err := chromedp.Run(tabCtx, actions...)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	return resp
}

func (resp *Response) Render() *Response {
	if resp == nil {
		return resp
	}
	return resp.CustomRender(nil, nil, nil)
}

// CustomRender 支持各类Action接口实现
func (resp *Response) CustomRender(targetListenerCallbacks []func(ev interface{}), flags []chromedp.ExecAllocatorOption, actions ...chromedp.Action) *Response {
	if resp == nil {
		return resp
	}
	if actions == nil {
		actions = []chromedp.Action{chromedp.OuterHTML("html", &resp.Html, chromedp.ByQuery)}
	} else {
		actions = append(actions, chromedp.OuterHTML("html", &resp.Html, chromedp.ByQuery))
	}
	return resp.render(targetListenerCallbacks, flags, actions...)
}

// Snapshot quality: false->jpeg | true->png
func (resp *Response) Snapshot(fullscreen bool, png bool) *[]byte {
	var buf = new([]byte)
	quality := int(*(*int8)(unsafe.Pointer(&png))) * 100
	var screenShotAction chromedp.Action
	if fullscreen {
		screenShotAction = chromedp.FullScreenshot(buf, quality)
	} else {
		screenShotAction = chromedp.CaptureScreenshot(buf)
	}
	resp.render(nil, nil, chromedp.Tasks{
		screenShotAction,
	})
	return buf
}

func actionFuncSetCookies(resp *Response) chromedp.ActionFunc {
	domain := &resp.Request.URL.Host
	cookies := &resp.Session.Cookies
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
}

func actionBypassWebDriver() chromedp.ActionFunc {
	return func(cxt context.Context) error {
		_, err := page.AddScriptToEvaluateOnNewDocument("Object.defineProperty(navigator, 'webdriver', { get: () => false, });").Do(cxt)
		if err != nil {
			return err
		}
		return nil
	}
}
