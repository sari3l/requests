package requests

import (
	"github.com/pkg/errors"
	"github.com/sari3l/requests/ext"
	"github.com/sari3l/requests/types"
	"io/ioutil"
	"log"
	"net/http"
	nUrl "net/url"
	"strconv"
	"strings"
)

// 后续转为接口，需要引入其他adapter
type adapter struct {
}

func (a *adapter) send(client *http.Client, prep *prepareRequest, hooks types.HooksDict) *Response {
	req := &Request{Request: &http.Request{Proto: "HTTP/1.1"}}

	req.Method = prep.method

	url, _ := nUrl.Parse(prep.url)
	req.URL = url

	if prep.headers != nil {
		req.Header = *prep.headers
		if host := req.Header.Get("Host"); host != "" {
			req.Host = host
		}
		if lens := req.Header.Get("Content-Length"); lens != "" {
			length, err := strconv.ParseInt(lens, 10, 64)
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
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}

	response := a.buildResponse(req.Request, resp)
	if response == nil {
		return nil
	}

	responseHandle := ext.DisPatchHook("response", hooks, *response).(Response)

	return &responseHandle
}

func (a *adapter) buildResponse(req *http.Request, resp *http.Response) *Response {

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}
	defer resp.Body.Close()

	encoding := resp.Header.Get("Content-Encoding")
	if err = decompressRaw(&raw, encoding); err != nil {
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}

	r := &Response{
		Ok:       resp.StatusCode == 200,
		Response: resp,
		Raw:      raw,
		Html:     string(raw),
		cookies:  append(resp.Cookies(), req.Cookies()...),
	}

	return r
}
