package ext

import (
	"github.com/sari3l/requests/types"
	"io"
)

func AllowRedirects(allowRedirects bool) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.AllowRedirects = allowRedirects
	}
}

func Auth(auth types.AuthInter) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Auth = auth
	}
}

func Cookies(cookies types.Dict) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Cookies = cookies
	}
}

func Data(data types.Dict) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Data = data
	}
}

func Files(files types.Dict) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Files = files
	}
}

func Headers(headers types.Dict) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Headers = headers
	}
}

func Hooks(hooksDict types.HooksDict) types.Ext {
	if len(hooksDict) == 0 {
		hooksDict = DefaultHooks()
	}
	return func(ep *types.ExtensionPackage) {
		ep.Hooks = hooksDict
	}
}

func Json(json map[string]any) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Json = json
	}
}

func Params(params types.Dict) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Params = params
	}
}

func Proxy(proxy string) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Proxy = proxy
	}
}

func Stream(stream io.Reader) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Stream = stream
	}
}

func Timeout(timeout int) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Timeout = timeout
	}
}

func Verify(verify bool) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Verify = verify

	}
}
