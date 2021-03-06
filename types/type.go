package types

import (
	"io"
)

type Dict map[string]string
type List []string
type Json map[string]any

type Hook func(object any) (error, any)
type HooksDict map[string][]Hook

type AuthInter interface {
	Format(p any) error
}

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
}
