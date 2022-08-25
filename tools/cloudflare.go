package tools

import (
	"github.com/sari3l/requests"
	"github.com/sari3l/requests/types"
	"net/url"
	"strings"
)

// HookCloudFlareWorkFunc 利用work反代进行访问，主要设置以下两个header
// Px-Host 会自从从原始请求中提取目标域名，也可通过显式设置字典值进行覆盖
// Px-IP 伪装IP，合法即可
// 使用此 hook 最重要保证链接到 cloudflare 的出口是白的
// 参考：https://github.com/jychp/cloudflare-bypass

func HookCloudFlareWorkerFunc(workerHost string, headers types.Dict) types.Hook {
	return func(request any) (error, any) {
		proxy, _ := url.Parse(workerHost)
		reqUrl := request.(requests.Request).URL
		if headers != nil {
			if headers["Px-Host"] == "" {
				headers["Px-Host"] = reqUrl.Host
			}
			if headers["Px-IP"] == "" {
				headers["Px-IP"] = RandomIPv4()
			}
		}
		request.(requests.Request).URL.Scheme = "https"
		if proxy.Host != "" {
			request.(requests.Request).URL.Host = proxy.Host
		} else {
			// 冗余处理，当url解析不带协议头的URL，会存在Path中，需要切片取host部分
			request.(requests.Request).URL.Host = strings.Split(proxy.Path, "/")[0]
		}
		reqHeader := request.(requests.Request).Header
		for key, value := range headers {
			reqHeader.Set(key, value)
		}
		return nil, nil
	}
}
