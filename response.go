package requests

import (
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
	Content string
	History []*Response
	Time    int64
}

func (resp *Response) Json() gjson.Result {
	if resp == nil {
		return gjson.Result{}
	}
	return gjson.Parse(resp.Content)
}

func (resp *Response) Text() string {
	text := ""

	domDoc := html.NewTokenizer(strings.NewReader(resp.Content))
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
	links := linkRegexCompiled.FindAllString(resp.Content, -1)
	originUrl := resp.Request.URL
	return *processLinks(originUrl, &links)
}
