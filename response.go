package requests

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/sari3l/requests/parser"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	*http.Response
	cookies []*http.Cookie
	Ok      bool
	Raw     []byte
	Html    string
	History []*Response
	Time    int64
	context *context.Context
	closer  *func()
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

// chromeless 页面渲染

func (resp *Response) Render() error {
	if resp.context == nil {
		ctx, cancel := chromedp.NewContext(
			context.Background(),
		)
		resp.context = &ctx
		resp.closer = (*func())(&cancel)
	}
	err := chromedp.Run(*resp.context,
		chromedp.Navigate(resp.Request.URL.String()),
		chromedp.OuterHTML("html", &resp.Html, chromedp.ByQuery),
	)
	if err != nil {
		return err
	}
	return nil
}

func (resp *Response) Close() {
	if resp.closer != nil {
		(*resp.closer)()
	}
}
