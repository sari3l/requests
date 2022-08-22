package requests

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/sari3l/requests/parser"
	"github.com/sari3l/requests/types"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	*http.Response
	Session *session
	cookies []*http.Cookie
	Ok      bool
	Raw     []byte
	Html    string
	History []*Response
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

func (resp *Response) URLs() []string {
	links := linkRegexCompiled.FindAllString(resp.Html, -1)
	originUrl := resp.Request.URL
	return *processLinks(originUrl, &links)
}

// Render chromeless 页面渲染，后续考虑将 chromedp 作为tools接入了，提供更多可选项
func (resp *Response) Render(useExtCookies bool) error {
	var flags = make([]chromedp.ExecAllocatorOption, 0)
	flags = append(flags,
		chromedp.Flag("ignore-certificate-errors", !resp.Session.Verify),
		chromedp.Flag("headless", true),
		chromedp.ProxyServer(resp.Session.Proxy),
		chromedp.UserAgent(resp.Request.UserAgent()),
	)

	setCookiesAction := actionFuncSetCookies(useExtCookies, resp.Request.URL.Host, &resp.Session.Cookies, &resp.cookies)

	opts := append(chromedp.DefaultExecAllocatorOptions[:], flags...)
	chromeCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(chromeCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	err := chromedp.Run(ctx,
		setCookiesAction,
		chromedp.Navigate(resp.Request.URL.String()),
		// 这里添加睡眠时间
		//chromedp.Sleep(10*time.Second),
		chromedp.OuterHTML("html", &resp.Html, chromedp.ByQuery),
	)
	if err != nil {
		return err
	}
	return nil
}

func actionFuncSetCookies(useExtCookies bool, domain string, cookiesDict *types.Dict, cookiesResp *[]*http.Cookie) chromedp.ActionFunc {
	if useExtCookies || cookiesResp == nil || len(*cookiesResp) == 0 {
		return func(ctx context.Context) error {
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			for name, value := range *cookiesDict {
				err := network.SetCookie(name, value).
					WithExpires(&expr).
					WithDomain(domain).
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
	} else {
		return func(ctx context.Context) error {
			for _, cookie := range *cookiesResp {
				err := network.SetCookie(cookie.Name, cookie.String()).
					WithExpires((*cdp.TimeSinceEpoch)(&cookie.Expires)).
					WithDomain(cookie.Domain).
					WithPath(cookie.Path).
					WithHTTPOnly(cookie.HttpOnly).
					WithSecure(cookie.Secure).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}
}
