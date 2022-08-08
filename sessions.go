package requests

import (
	"crypto/tls"
	"github.com/jinzhu/copier"
	"github.com/sari3l/requests/ext"
	"github.com/sari3l/requests/types"
	"log"
	"net/http"
	nUrl "net/url"
	"time"
)

type session struct {
	types.ExtensionPackage
	Client       *http.Client
	adapter      *adapter
	cacheRequest *prepareRequest
}

func Session(timeout int, proxy string, redirect bool, verify bool) *session {
	s := &session{
		adapter: &adapter{},
		Client:  &http.Client{},
	}

	s.Timeout = timeout
	s.Proxy = proxy
	s.AllowRedirects = redirect
	s.Verify = verify

	return s
}

func (s *session) init(method string, url string, exts *[]types.Ext) *session {
	s.Method = method
	s.Url = url

	for _, fn := range *exts {
		fn(&s.ExtensionPackage)
	}

	return s
}

func (s *session) request() *Response {
	err, prep := s.prepareRequest()
	if err != nil {
		log.Println(err)
		return nil
	}

	if err = s.prepareClient(); err != nil {
		log.Println(err)
		return nil
	}

	return s.Send(prep)
}

func (s *session) prepareRequest() (error, *prepareRequest) {
	err, prep := PrepareRequest(s.Method, s.Url, s.Params, s.Headers, s.Cookies, s.Data, s.Json, s.Files, s.Stream, s.Auth, s.Hooks)
	if err != nil {
		return err, nil
	}

	return nil, prep
}

func (s *session) prepareClient() error {
	_ = s.prepareTimeout()
	if err := s.prepareProxy(); err != nil {
		return err
	}
	_ = s.prepareRedirect()
	_ = s.prepareVerify()

	return nil
}

func (s *session) Send(prep *prepareRequest) *Response {
	var err error

	// 计时开机
	startTime := time.Now().UnixMilli()
	// 后续根据协议，切换adapter
	err, r := s.adapter.send(s.Client, prep, s.Hooks)
	if err != nil {
		log.Println(err)
		return nil
	}
	usedTime := time.Now().UnixMilli() - startTime

	r.Time = usedTime

	history := make([]*Response, 0)
	if s.AllowRedirects {
		s.cacheRequest = prep
		err, history = s.resolveRedirects(r)
	}

	if len(history) > 0 {
		history = append([]*Response{r}, history...)
		r = history[len(history)-1]
		r.History = history[:len(history)-1]
	}

	return r
}

func (s *session) RegisterHook(key string, hook types.Hook) error {
	if s.Hooks == nil {
		s.Hooks = ext.DefaultHooks()
	}
	return ext.RegisterHook(&s.Hooks, key, hook)
}

func (s *session) resolveRedirects(resp *Response) (error, []*Response) {
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
			return err, history
		}
		redirectPrep.cookies = resp.Cookies()
		s.AllowRedirects = false
		resp = s.Send(redirectPrep)
		url = resp.Header.Get("Location")
		history = append(history, resp)
	}
	return nil, history
}

func (s *session) prepareRedirect() error {
	s.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return nil
}

func (s *session) prepareProxy() error {
	if s.Proxy == "" {
		return nil
	}
	_proxy, err := nUrl.Parse(s.Proxy)
	if err != nil {
		return err
	}
	if s.Client.Transport != nil {
		s.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(_proxy)
	} else {
		s.Client.Transport = &http.Transport{
			Proxy: http.ProxyURL(_proxy),
		}
	}
	return nil
}

func (s *session) prepareTimeout() error {
	if s.Timeout == 0 {
		s.Timeout = DefaultTimeout
	}
	s.Client.Timeout = time.Duration(s.Timeout) * time.Second
	return nil
}

func (s *session) prepareVerify() error {
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

func (s *session) SetVerify(verify bool) error {
	s.Verify = verify
	return s.prepareVerify()
}

func (s *session) Get(url string, ext ...types.Ext) *Response {
	return s.init("Get", url, &ext).request()
}

func (s *session) Post(url string, ext ...types.Ext) *Response {
	return s.init("Post", url, &ext).request()
}

func (s *session) Put(url string, ext ...types.Ext) *Response {
	return s.init("Put", url, &ext).request()
}

func (s *session) Delete(url string, ext ...types.Ext) *Response {
	return s.init("Delete", url, &ext).request()
}

func (s *session) Head(url string, ext ...types.Ext) *Response {
	return s.init("Head", url, &ext).request()
}

func (s *session) Options(url string, ext ...types.Ext) *Response {
	return s.init("Option", url, &ext).request()
}
