package requests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sari3l/requests/ext"
	"github.com/sari3l/requests/internal/decoder"
	"github.com/sari3l/requests/internal/processbar"
	"github.com/sari3l/requests/internal/tracer"
	"github.com/sari3l/requests/types"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	nUrl "net/url"
	"strconv"
	"strings"
	"time"
)

// 后续转为接口，需要引入其他adapter
type adapter struct {
	context      context.Context
	tracerEnable bool
}

func (a *adapter) send(client *http.Client, prep *prepareRequest, hooks *types.HooksDict) *Response {
	trace := prepareTracer(a.tracerEnable, fmt.Sprintf("%s >>>>> %s", prep.method, prep.url))
	req := &Request{Request: &http.Request{Proto: prep.proto}}
	a.context = httptrace.WithClientTrace(a.context, trace.ClientTrace)
	req.Request = req.WithContext(a.context)
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

	requestHandle := ext.DisPatchHook("request", *hooks, *req).(Request)
	clientHandle := ext.DisPatchHook("client", *hooks, *client).(http.Client)

	if a.tracerEnable {
		trace.StartTime = time.Now()
	}

	resp, err := clientHandle.Do(requestHandle.Request)

	if a.tracerEnable {
		trace.EndTime = time.Now()
		log.Println(trace.ToString())
	}

	if resp == nil || err != nil {
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}

	response := a.buildResponse(req.Request, resp)
	if response == nil {
		return nil
	}

	responseHandle := ext.DisPatchHook("response", *hooks, *response).(Response)

	return &responseHandle
}

func (a *adapter) buildResponse(req *http.Request, resp *http.Response) *Response {
	// 允许接入多个writer
	buf := &bytes.Buffer{}
	opts, ok := a.context.Value("processOptions").([]processbar.Option)
	if !ok {
		opts = []processbar.Option{}
	}
	pb := processbar.NewProcessBar(resp.ContentLength, opts...)
	_, err := io.Copy(io.MultiWriter(buf, pb), resp.Body)

	if err != nil {
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}
	defer resp.Body.Close()
	raw := buf.Bytes()
	//raw, err := ioutil.ReadAll(resp.Body)

	// 判断raw长度是否需要解码
	encoding := resp.Header.Get("Content-Encoding")
	if err = decoder.DecompressRaw(&raw, encoding); err != nil {
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}

	r := &Response{
		Ok:       resp.StatusCode == 200,
		Response: resp,
		Raw:      &raw,
		Html:     string(raw),
		cookies:  append(resp.Cookies(), req.Cookies()...),
	}

	return r
}

func prepareTracer(enable bool, prefix string) *tracer.Tracer {
	newTracer := tracer.Tracer{Enable: enable}
	newTracer.Output.Prefix = prefix
	trace := newTracer.Init()
	return trace
}
