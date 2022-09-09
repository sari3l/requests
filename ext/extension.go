package ext

import (
	"github.com/sari3l/requests/internal/processbar"
	"github.com/sari3l/requests/types"
	"io"
)

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

func CipherSuites(cipherSuites []uint16) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.CipherSuites = cipherSuites
	}
}

func Form(form types.Dict) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Form = form
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

func Protocol(protocol string) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Protocol = protocol
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

func ProcessOptions(opts ...processbar.Option) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.ProcessOptions = opts
	}
}

func Proxy(proxy string) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Proxy = proxy
	}
}

func Redirect(allow bool) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Redirect = allow
	}
}

func Stream(stream io.Reader) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Stream = stream
	}
}

func Timeout(second int) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Timeout = second
	}
}

func Tracer(enable bool) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Tracer = enable
	}
}

func Verify(verify bool) types.Ext {
	return func(ep *types.ExtensionPackage) {
		ep.Verify = verify
	}
}
