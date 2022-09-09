package requests

import (
	"bytes"
	eJson "encoding/json"
	"errors"
	"github.com/sari3l/requests/types"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	nUrl "net/url"
	"os"
	"strconv"
	"strings"
)

type prepareRequest struct {
	method  string
	url     string
	proto   string
	headers *http.Header
	cookies []*http.Cookie
	body    *io.ReadCloser
	hooks   types.HooksDict
}

func PrepareRequest(proto string, method string, url string, params types.Dict, headers types.Dict, cookies types.Dict, form types.Dict, json map[string]any, files types.Dict, stream io.Reader, auth types.AuthInter) (error, *prepareRequest) {
	var err error
	_prepareRequest := new(prepareRequest)
	_prepareRequest.prepareProto(proto)
	_prepareRequest.prepareMethod(method)
	if err = _prepareRequest.prepareUrl(url, params); err != nil {
		return err, nil
	}
	if err = _prepareRequest.prepareBody(form, files, json, stream); err != nil {
		return err, nil
	}
	if err = _prepareRequest.prepareHeaders(headers); err != nil {
		return err, nil
	}
	if err = _prepareRequest.prepareCookies(cookies); err != nil {
		return err, nil
	}
	if err = _prepareRequest.prepareAuth(auth); err != nil {
		return err, nil
	}

	if err != nil {
		return err, nil
	}

	return nil, _prepareRequest
}

func (prep *prepareRequest) prepareProto(proto string) {
	if proto != "" {
		prep.proto = proto
	} else if strings.HasSuffix(proto, "2") {
		prep.proto = "HTTP/2"
	} else {
		prep.proto = "HTTP/1.1"
	}
}

func (prep *prepareRequest) prepareMethod(method string) {
	prep.method = strings.ToUpper(method)
}

func (prep *prepareRequest) prepareUrl(urlRaw string, params types.Dict) error {
	urlRaw = strings.Trim(urlRaw, " ")
	if strings.Contains(urlRaw, ":") && strings.ToLower(urlRaw[:4]) != "http" {
		prep.url = urlRaw
		return nil
	}

	url, err := nUrl.Parse(urlRaw)
	if err != nil {
		return err
	}
	if url.Scheme == "" {
		return errors.New("无协议头")
	}

	if url.Host == "" {
		return errors.New("无有效域名")
	} else if strings.HasPrefix(url.Host, "*") || strings.HasPrefix(url.Host, ".") {
		return errors.New("错误域名")
	}

	if url.Path == "" {
		url.Path = "/"
	}
	if params != nil {
		query := url.Query()
		for k, v := range params {
			query.Add(k, v)
		}
		url.RawQuery = query.Encode()
	}

	prep.url = url.String()
	return nil
}

// 已知问题，不按初始化排序，google的高傲尽显，So fuck you google 🖕
// https://go-review.googlesource.com/c/go/+/105755
func (prep *prepareRequest) prepareHeaders(headers types.Dict) error {
	if headers != nil {
		if prep.headers == nil {
			prep.headers = &http.Header{}
		}
		for k, v := range headers {
			// 在 go/src/net/http/transport.go:504 -> roundTrip 有header头检查
			prep.headers.Set(k, v)
		}
	} else {
		prep.headers = defaultHeaders()
	}
	return nil
}

func (prep *prepareRequest) prepareCookies(cookies types.Dict) error {
	if cookies == nil {
		return nil
	}
	for k, v := range cookies {
		cookie := &http.Cookie{Name: k, Value: v, Path: "/", Domain: ""}
		prep.cookies = append(prep.cookies, cookie)
	}
	tmpReq := http.Request{Header: *prep.headers}
	for _, c := range prep.cookies {
		tmpReq.AddCookie(c)
	}
	prep.headers.Set("Cookie", tmpReq.Header.Get("Cookie"))
	return nil
}

func (prep *prepareRequest) prepareBody(form, files types.Dict, json map[string]any, stream io.Reader) error {
	var closer io.ReadCloser
	contentType := ""
	contentLength := 0
	// json
	if form == nil && json != nil {
		jsonByte, err := eJson.Marshal(json)
		if err != nil {
			return err
		}
		closer = io.NopCloser(bytes.NewReader(jsonByte))
		contentLength = len(jsonByte)
		contentType = "application/json"
	}

	if stream != nil {
		closer = io.NopCloser(stream)
		contentType = "application/octet-stream"
	} else if files != nil {
		buffer := bytes.Buffer{}
		writer := multipart.NewWriter(&buffer)
		for field, filename := range files {
			part, err := writer.CreateFormFile(field, filename)
			if err != nil {
				panic(err)
			}
			file, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(part, file)
			if err != nil {
				panic(err)
			}
		}
		if form != nil {
			for k, v := range form {
				_ = writer.WriteField(k, v)
			}
		}
		defer writer.Close()
		closer = ioutil.NopCloser(bytes.NewReader(buffer.Bytes()))
		contentLength = buffer.Len()
		contentType = "multipart/form-form"
	} else if form != nil {
		dataValues := nUrl.Values{}
		for k, v := range form {
			dataValues.Set(k, v)
		}
		_dataEncoded := dataValues.Encode()
		closer = io.NopCloser(strings.NewReader(_dataEncoded))
		contentLength = len(_dataEncoded)
		contentType = "application/x-www-form-urlencoded"
	}

	if contentType != "" {
		prep.headers.Set("Content-Type", contentType)
	}
	if contentLength != 0 {
		prep.headers.Set("Content-Length", strconv.Itoa(contentLength))
	}
	prep.body = &closer
	return nil
}

func (prep *prepareRequest) prepareAuth(auth types.AuthInter) error {
	if auth == nil {
		return nil
	}
	return auth.Format(prep.headers)
}
