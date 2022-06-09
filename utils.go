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
	"regexp"
	"strings"
)

const Version = "v1.0.0"

const DefaultTimeout = 5 // time.Second

var cleanHeaderRegexStr = regexp.MustCompile(`^\S[^\r\n]*$|^$`)

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
