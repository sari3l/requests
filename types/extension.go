package types

import "io"

type Ext func(ep *ExtensionPackage)
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
	Render         bool
}
