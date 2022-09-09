package types

import (
	"github.com/sari3l/requests/internal/processbar"
	"io"
)

type Ext func(ep *ExtensionPackage)
type ExtensionPackage struct {
	Method         string    // 隐性参数
	Url            string    // 直传参数
	Auth           AuthInter // 以下均为扩展参数
	CipherSuites   []uint16
	Cookies        Dict
	Form           Dict
	Files          Dict
	Headers        Dict
	Hooks          HooksDict
	Json           map[string]any
	Protocol       string
	LimitLength    int64
	Params         Dict
	ProcessOptions []processbar.Option
	Proxy          string
	Redirect       bool
	Stream         io.Reader
	Timeout        int
	Tracer         bool
	Verify         bool
}
