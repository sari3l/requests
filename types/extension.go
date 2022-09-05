package types

import "io"

type Ext func(ep *ExtensionPackage)
type ExtensionPackage struct {
	Method         string // 隐性参数
	Url            string // 直传参数
	AllowRedirects bool   // 以下均为扩展参数
	Auth           AuthInter
	CipherSuites   []uint16
	Cookies        Dict
	Data           Dict
	Files          Dict
	Headers        Dict
	Hooks          HooksDict
	Proto          string
	Json           map[string]any
	Params         Dict
	Proxy          string
	Stream         io.Reader
	Timeout        int
	Verify         bool
}
