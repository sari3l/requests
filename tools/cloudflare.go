package tools

import (
	"github.com/sari3l/requests"
	"github.com/sari3l/requests/types"
	"net/url"
	"strings"
)

// HookCloudFlareWorkFunc 利用work反代进行访问，主要设置以下两个header
// 使用此 hook 最重要保证链接到 cloudflare 的出口是白的
// 参考：https://github.com/jychp/cloudflare-bypass

func HookCloudFlareWorkerFunc(workerHost string, headers types.Dict) types.Hook {
	return func(request any) (error, any) {
		worker, _ := url.Parse(workerHost)
		var host string
		request.(requests.Request).URL.Scheme = "https"
		if worker.Host != "" {
			host = worker.Host
		} else {
			// 冗余处理，当url解析不带协议头的URL，会存在Path中，需要切片取host部分
			host = strings.Split(worker.Path, "/")[0]
		}
		request.(requests.Request).Request.Host = host
		request.(requests.Request).URL.Host = host
		reqHeader := request.(requests.Request).Header
		for key, value := range headers {
			reqHeader.Set(key, value)
		}
		return nil, nil
	}
}
