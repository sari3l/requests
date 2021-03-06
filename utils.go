package requests

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	"io/ioutil"
	"net/http"
	nUrl "net/url"
	"regexp"
	"strings"
)

const Version = "v1.1.5.1"

const DefaultTimeout = 5 // time.Second

var cleanHeaderRegexStr = regexp.MustCompile(`^\S[^\r\n]*$|^$`)
var linkRegexCompiled = regexp.MustCompile(`(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;|*()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|/][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:\?[^"|']{0,}|)))(?:"|')`)

func checkHeaderValidity(key, value string) error {
	if !cleanHeaderRegexStr.MatchString(value) {
		return errors.New(fmt.Sprintf("header %s 值结尾错误", key))
	}
	return nil
}

func defaultHeaders() *http.Header {
	headers := &http.Header{}
	headers.Add("User-Agent", fmt.Sprintf("sari3l/requests %s", Version))
	return headers
}

func decompressRaw(raw *[]byte, encoding string) error {
	if encoding == "" {
		return nil
	}
	switch strings.ToLower(encoding) {
	case "gzip":
		return decompressGzip(raw)
	case "deflate":
		return decompressDeflate(raw)
	case "br":
		return decompressBr(raw)
	}
	return nil
}

func decompressGzip(raw *[]byte) error {
	reader, err := gzip.NewReader(bytes.NewReader(*raw))
	if err != nil {
		return err
	}
	defer reader.Close()
	*raw, err = ioutil.ReadAll(reader)
	return err
}

func decompressDeflate(raw *[]byte) error {
	var err error
	reader := flate.NewReader(bytes.NewReader(*raw))
	defer reader.Close()
	*raw, err = ioutil.ReadAll(reader)
	return err
}

func decompressBr(raw *[]byte) error {
	var err error
	r := brotli.NewReader(bytes.NewReader(*raw))
	*raw, err = ioutil.ReadAll(r)
	return err
}

func processLinks(url *nUrl.URL, links *[]string) *[]string {
	for index, link := range *links {
		link = strings.Trim(link, "\"")
		link = strings.Trim(link, "'")
		if len(link) >= 2 && link[0:2] == "//" {
			(*links)[index] = url.Scheme + ":" + link
		} else if len(link) >= 4 && link[0:4] == "http" {
			continue
		} else if len(link) >= 2 && link[0:2] != "//" {
			if link[0:1] == "/" {
				(*links)[index] = url.Scheme + "://" + url.Host + link
			} else if link[0:1] == "." {
				if link[0:2] == ".." {
					(*links)[index] = url.Scheme + "://" + url.Host + link[2:]
				} else {
					(*links)[index] = url.Scheme + "://" + url.Host + link[1:]
				}
			} else {
				(*links)[index] = url.Scheme + "://" + url.Host + "/" + link
			}
		} else {
			continue
		}
	}
	return links
}
