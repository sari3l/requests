package decoder

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"github.com/andybalholm/brotli"
	"io/ioutil"
	"reflect"
	"strings"
)

// gzip https://en.wikipedia.org/wiki/Gzip
var gzipFlag = []byte{0x1f, 0x8b}

// DecompressRaw 优化自动识别，目前有点丑陋，可以封装成一个接口向外提供
// 还需要判断需不需要转码，目前没木有
func DecompressRaw(raw *[]byte, encoding string) error {
	if raw == nil || len(*raw) == 0 {
		return nil
	}
	if encoding == "" {
		// 解析压缩魔术头
		if reflect.DeepEqual((*raw)[:2], gzipFlag) {
			encoding = "gzip"
		} else {
			return nil
		}
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
