package requests

import (
	"github.com/sari3l/requests/internal/parser"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type Response struct {
	*http.Response
	Session *Session
	cookies []*http.Cookie
	History []*Response
	Html    string
	Ok      bool
	Raw     *[]byte
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
