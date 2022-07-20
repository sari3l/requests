package requests

import (
	"fmt"
	"github.com/sari3l/requests/ext"
	"github.com/sari3l/requests/types"
	"io/ioutil"
	"net/http"
	nUrl "net/url"
	"strconv"
	"strings"
)

// 后续转为接口，需要引入其他adapter
type adapter struct {
}

func (a *adapter) send(client *http.Client, prep *prepareRequest, hooks types.HooksDict) (error, *Response) {
	req := &Request{Request: &http.Request{Proto: "HTTP/1.1"}}

	req.Method = prep.method

	url, _ := nUrl.Parse(prep.url)
	req.URL = url

	if prep.headers != nil {
		req.Header = *prep.headers
		if req.Header.Get("Content-Length") != "" {
			length, err := strconv.ParseInt(req.Header.Get("Content-Length"), 10, 64)
			if err == nil {
				req.ContentLength = length
			}
		}
		if strings.ToLower(req.Header.Get("Transfer-Encoding")) == "chunked" {
			req.ContentLength = 0
		}
	}

	if prep.body != nil {
		req.Body = *prep.body
	}

	requestHandle := ext.DisPatchHook("request", hooks, *req).(Request)
	clientHandle := ext.DisPatchHook("client", hooks, *client).(http.Client)

	resp, err := clientHandle.Do(requestHandle.Request)
	if resp == nil || err != nil {
		return err, nil
	}

	err, response := a.buildResponse(req.Request, resp)
	if err != nil {
		fmt.Printf("%v", err)
		return err, nil
	}

	responseHandle := ext.DisPatchHook("response", hooks, *response).(Response)

	return nil, &responseHandle
}

func (a *adapter) buildResponse(req *http.Request, resp *http.Response) (error, *Response) {

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	encoding := resp.Header.Get("Content-Encoding")
	_ = decompressRaw(&raw, encoding)

	r := &Response{
		Ok:       resp.StatusCode == 200,
		Response: resp,
		Raw:      raw,
		Content:  string(raw),
		cookies:  append(resp.Cookies(), req.Cookies()...),
	}

	return nil, r
}
