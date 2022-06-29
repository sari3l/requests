package ext

import (
	"io"
)

type Ext func(ep *ExtensionPackage)
type Dict map[string]string
type List []string

type ExtensionPackage struct {
	Method         string // 隐性参数
	Url            string // 直传参数
	Data           Dict   // 以下均为扩展参数
	Params         Dict
	Proxy          string
	Timeout        int
	Headers        Dict
	Cookies        Dict
	AllowRedirects bool
	Json           map[string]any
	Files          Dict
	Stream         io.Reader
	Auth           AuthInter
	Hooks          HooksDict
	Verify         bool
}

func AllowRedirects(allowRedirects bool) Ext {
	return func(ep *ExtensionPackage) {
		ep.AllowRedirects = allowRedirects
	}
}

func Auth(auth AuthInter) Ext {
	return func(ep *ExtensionPackage) {
		ep.Auth = auth
	}
}

func Cookies(cookies Dict) Ext {
	return func(ep *ExtensionPackage) {
		ep.Cookies = cookies
	}
}

func Data(data Dict) Ext {
	return func(ep *ExtensionPackage) {
		ep.Data = data
	}
}

func Files(files Dict) Ext {
	return func(ep *ExtensionPackage) {
		ep.Files = files
	}
}

func Headers(headers Dict) Ext {
	return func(ep *ExtensionPackage) {
		ep.Headers = headers
	}
}

func Hooks(hooksDict HooksDict) Ext {
	if len(hooksDict) == 0 {
		hooksDict = DefaultHooks()
	}
	return func(ep *ExtensionPackage) {
		ep.Hooks = hooksDict
	}
}

func Json(json map[string]any) Ext {
	return func(ep *ExtensionPackage) {
		ep.Json = json
	}
}

func Params(params Dict) Ext {
	return func(ep *ExtensionPackage) {
		ep.Params = params
	}
}

func Proxy(proxy string) Ext {
	return func(ep *ExtensionPackage) {
		ep.Proxy = proxy
	}
}

func Stream(stream io.Reader) Ext {
	return func(ep *ExtensionPackage) {
		ep.Stream = stream
	}
}

func Timeout(timeout int) Ext {
	return func(ep *ExtensionPackage) {
		ep.Timeout = timeout
	}
}

func Verify(verify bool) Ext {
	return func(ep *ExtensionPackage) {
		ep.Verify = verify

	}
}
