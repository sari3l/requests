package requests

import (
	"context"
	"crypto/tls"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/sari3l/requests/ext"
	"github.com/sari3l/requests/internal/processbar"
	tracer2 "github.com/sari3l/requests/internal/tracer"
	"github.com/sari3l/requests/types"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"net/http/httptrace"
	nUrl "net/url"
	"time"
)

type Session struct {
	types.ExtensionPackage
	Client       *http.Client
	adapter      *adapter
	cacheRequest *prepareRequest
	tracer       *tracer2.Tracer
}

func HTMLSession() *Session {
	s := &Session{
		adapter: &adapter{context: context.Background()},
		Client: &http.Client{
			Transport: new(http.Transport),
		},
	}
	s.Timeout = 5
	s.Redirect = true
	s.Verify = true
	return s
}

func (s *Session) init(method string, url string, exts *[]types.Ext) *Session {
	s.Method = method
	s.Url = url

	for _, fn := range *exts {
		fn(&s.ExtensionPackage)
	}

	return s
}

func (s *Session) request() *Response {
	var err error
	err, s.cacheRequest = s.prepareRequest()
	if err != nil {
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}

	if err = s.prepareClient(); err != nil {
		log.Printf("%+v\n", errors.WithStack(err))
		return nil
	}

	s.prepareTracer()
	s.prepareProcessOptions()
	s.prepareContext()

	return s.Send(s.cacheRequest)
}

func (s *Session) prepareRequest() (error, *prepareRequest) {
	err, prep := PrepareRequest(s.Protocol, s.Method, s.Url, s.Params, s.Headers, s.Cookies, s.Form, s.Json, s.Files, s.Stream, s.Auth)
	if err != nil {
		return err, nil
	}

	return nil, prep
}

func (s *Session) prepareClient() error {
	_ = s.prepareTimeout()
	_ = s.prepareCipherSuites()
	if err := s.prepareProxy(); err != nil {
		return err
	}
	_ = s.prepareRedirect()
	_ = s.prepareVerify()
	return nil
}

func (s *Session) Send(prep *prepareRequest) *Response {
	// 计时开机
	startTime := time.Now()
	// 后续根据协议，切换 adapter（实际go对应client配置）
	resp := s.adapter.send(s.Client, prep, &s.Hooks)
	history := make([]*Response, 0)

	if resp != nil {
		if s.Redirect {
			s.cacheRequest = prep
			history = s.resolveRedirects(resp)
		}

		if len(history) > 0 {
			history = append([]*Response{resp}, history...)
			resp = history[len(history)-1]
			resp.History = history[:len(history)-1]
		}
	}

	endTime := time.Now()
	if s.tracer != nil {
		s.tracer.StartTime = startTime
		s.tracer.EndTime = endTime
	}

	if s.Tracer == true {
		log.Println(s.tracer.ToString())
	}

	if resp == nil {
		return nil
	}

	resp.Session = s
	return resp
}

func (s *Session) RegisterHook(key string, hook types.Hook) error {
	if s.Hooks == nil {
		s.Hooks = ext.DefaultHooks()
	}
	return ext.RegisterHook(&s.Hooks, key, hook)
}

func (s *Session) resolveRedirects(resp *Response) []*Response {
	var err error
	history := make([]*Response, 0)
	url := resp.Header.Get("Location")
	for url != "" {
		if u := string(url[0]); u == "/" || u == "." {
			uTmp := resp.Request.URL
			uTmp.Path = url
			url = uTmp.String()
		}
		redirectPrep := &prepareRequest{}
		_ = copier.CopyWithOption(&redirectPrep, &s.cacheRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if err = redirectPrep.prepareUrl(url, nil); err != nil {
			log.Printf("%+v\n", errors.WithStack(err))
			return history
		}
		redirectPrep.cookies = resp.Cookies()
		s.Redirect = false
		resp = s.Send(redirectPrep)
		if resp == nil {
			break
		}
		url = resp.Header.Get("Location")
		history = append(history, resp)
	}
	return history
}

func (s *Session) prepareRedirect() error {
	s.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return nil
}

// net/http/Transport 支持 http/https/socks5，如果使用其他协议则会转为 http
// proxy支持socks/socks5h，为了简便只使用其处理socks5h
func (s *Session) prepareProxy() error {
	if s.Proxy == "" {
		return nil
	}
	_proxy, err := nUrl.Parse(s.Proxy)
	if err != nil {
		return err
	}
	switch _proxy.Scheme {
	case "http", "https", "socks5":
		s.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(_proxy)
	case "socks5h":
		dialer, err := proxy.FromURL(_proxy, nil)
		if err != nil {
			return err
		}
		s.Client.Transport.(*http.Transport).DialContext = dialer.(proxy.ContextDialer).DialContext
	}
	// 设置 ProxyConnectHeader
	if ua := s.cacheRequest.headers.Get("User-Agent"); ua != "" {
		header := &http.Header{}
		header.Add("User-Agent", ua)
		s.Client.Transport.(*http.Transport).ProxyConnectHeader = *header
	}
	return nil
}

func (s *Session) prepareTimeout() error {
	if s.Timeout == 0 {
		s.Timeout = DefaultTimeout
	}
	s.Client.Timeout = time.Duration(s.Timeout) * time.Second
	return nil
}

func (s *Session) prepareVerify() error {
	if s.Verify == false {
		if s.Client.Transport != nil {
			s.Client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			s.Client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
	}
	return nil
}

func (s *Session) prepareTracer() {
	newTracer := tracer2.Tracer{Enable: s.Tracer}
	newTracer.Output.Url = s.Url
	newTracer.Output.Method = s.Method
	trace := newTracer.Init()
	s.tracer = trace
}

func (s *Session) prepareProcessOptions() {
	if s.ProcessOptions == nil {
		s.ProcessOptions = make([]processbar.Option, 0)
	}
	s.ProcessOptions = append(s.ProcessOptions, processbar.OptionsUrl(s.Url))
}

func (s *Session) prepareContext() {
	s.adapter.context = httptrace.WithClientTrace(s.adapter.context, s.tracer.ClientTrace)
	s.adapter.context = context.WithValue(s.adapter.context, "processOptions", s.ProcessOptions)
}

func (s *Session) prepareHooks(hooks types.HooksDict) {
	if s.Hooks == nil {
		s.Hooks = ext.DefaultHooks()
	}
	for event, _hooks := range hooks {
		if s.Hooks[event] != nil {
			s.Hooks[event] = append(s.Hooks[event], _hooks...)
		} else {
			s.Hooks[event] = _hooks
		}
	}
}

// crypto/tls 已做nil检查，直接传入即可
func (s *Session) prepareCipherSuites() error {
	s.Client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
		CipherSuites: s.CipherSuites,
	}
	return nil
}

func (s *Session) Get(url string, ext ...types.Ext) *Response {
	return s.init("Get", url, &ext).request()
}

func (s *Session) Post(url string, ext ...types.Ext) *Response {
	return s.init("Post", url, &ext).request()
}

func (s *Session) Put(url string, ext ...types.Ext) *Response {
	return s.init("Put", url, &ext).request()
}

func (s *Session) Delete(url string, ext ...types.Ext) *Response {
	return s.init("Delete", url, &ext).request()
}

func (s *Session) Head(url string, ext ...types.Ext) *Response {
	return s.init("Head", url, &ext).request()
}

func (s *Session) Options(url string, ext ...types.Ext) *Response {
	return s.init("Option", url, &ext).request()
}
