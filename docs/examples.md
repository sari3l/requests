# 示例

## Twitter API

```go
import (
    "fmt"
    "github.com/sari3l/requests"
    "github.com/sari3l/requests/ext"
    "github.com/sari3l/requests/types"
)

func main() {
    // Requests Bearer Token 
    auth := types.BasicAuth{Username: "o94KGT3MlbT...", Password: "fNbL2ukEGyvuGSM7bAuoq..."}
    data := types.Dict{
        "grant_type": "client_credentials",
    }
    resp := requests.Post("https://api.twitter.com/oauth2/token", ext.Auth(auth), ext.Data(data))

    // Requests with Twitter API 2.0
    if resp != nil && resp.Ok {
        fmt.Println(resp.Json())
        token := types.BearerAuth{Token: resp.Json().Get("access_token").Str}
        resp2 := requests.Get("https://api.twitter.com/2/users/by/username/Sariel_D", ext.Auth(token))
        fmt.Println(resp2.Json())
    }
}
```

## JA3 指纹

### 传统模式

通过生成`Session`、`PrepareRequest`初始化请求，并替换`Transport`内容实现

```go
import (
    "fmt"
    "github.com/CUCyber/ja3transport"
    "github.com/sari3l/requests"
)

func main() {
    session := requests.HTMLSession()
    _, prep := requests.PrepareRequest("HTTP/1.1", "get", "https://ja3er.com/json", nil, nil, nil, nil, nil, nil, nil, nil)
    tr, _ := ja3transport.NewTransport("771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53-10,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0")
    session.Client.Transport = tr
    resp := session.Send(prep)
    fmt.Print(resp.Html)
}
```

### Hook 模式 I

通过`Hook`方式替换`Transport`内容实现

```go
import (
    "fmt"
    "github.com/CUCyber/ja3transport"
    "github.com/sari3l/requests"
    "github.com/sari3l/requests/ext"
    "github.com/sari3l/requests/types"
    "net/http"
    "reflect"
)

func main() {
    hooks := types.HooksDict{
        "client": []types.Hook{modifyJa3Fingerprint},
    }

    resp := requests.Get("https://ja3er.com/json", ext.Hooks(hooks))
    fmt.Print(resp.Json())

}

func modifyJa3Fingerprint(client any) (error, any) {
    c := client.(http.Client)
    tr, _ := ja3transport.NewTransport("771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53-10,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0")
    reflect.ValueOf(&c).Elem().FieldByName("Transport").Set(reflect.ValueOf(tr))
    return nil, c
}
```

### Hook 模式 II

在`requests/tools`中提供了快速生成JA3Hook函数的方法，方便`Hook`使用

```go
import (
    "fmt"
    "github.com/sari3l/requests"
    "github.com/sari3l/requests/ext"
    "github.com/sari3l/requests/tools"
    "github.com/sari3l/requests/types"
)

func main() {
    ja3Hook := tools.HookClientJA3Func("771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53-10,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0")

    hooks := types.HooksDict{
        "client": []types.Hook{ja3Hook},
    }

    resp := requests.Get("https://ja3er.com/json", ext.Hooks(hooks))
    fmt.Print(resp.Json())
}
```

## CloudFlare

### Request with Workers

```go
import (
    "fmt"
    "github.com/sari3l/requests"
    "github.com/sari3l/requests/ext"
    "github.com/sari3l/requests/tools"
    "github.com/sari3l/requests/types"
)

func main() {
    headers := types.Dict {
        "Pragma": "no-cache",
        "Cache-Control": "no-cache",
        "Upgrade-Insecure-Requests": "1",
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
        "Accept-Encoding": "gzip, deflate",
        "Accept-Language": "zh-CN,zh;q=0.9",
        "Connection": "close",
    }
	
    hooks := types.HooksDict{
    "request": []types.Hook{tools.HookCloudFlareWorkFunc("https://delicate-xxx.sonymouse.workers.dev", types.Dict{
                "Px-Token": "mysecuretoken", // 自定义安全验证头
                "Px-Host": "www.google.com",
                "Px-IP": "1,2,3,4",
        })},
    }
	
    resp := requests.Get("https://www.google.com", ext.Headers(headers), ext.Hooks(hooks))
    fmt.Println(resp.Html)
}
```

## Crawler

### Image Downloader

#### 手动下载

```go
import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/chromedp/cdproto/input"
    "github.com/chromedp/cdproto/network"
    "github.com/chromedp/chromedp"
    "github.com/sari3l/requests"
    "github.com/sari3l/requests/ext"
    "io/ioutil"
    "strings"
    "time"
)

func main() {
    const baseUrl = "https://x.com/video?page=%d"
    for _, i := range []int{1, 2, 3} {
        fmt.Printf("正在爬取 %d 页\n", i)
        url := fmt.Sprintf(baseUrl, i)
        
        resp := requests.Get(url, ext.Timeout(5))
        resp.CustomRender(listenForNetworkEvent, nil,
            chromedp.EmulateViewport(1400, 2800, chromedp.EmulateScale(1)), 
            chromedp.Sleep(5*time.Second), 
            actionDispatchMouse())
    }
}

func actionDispatchMouse() chromedp.ActionFunc {
    return func(ctx context.Context) error {
        p := input.DispatchMouseEvent(input.MouseWheel, 200, 200)
        p = p.WithDeltaX(0)
        // 滚轮向下滚动1000单位
        p = p.WithDeltaY(float64(1000))
        return p.Do(ctx)
    }
}

type UrlResponse struct {
    Url string `json:"url"`
}

func listenForNetworkEvent(ev interface{}) {
    switch ev := ev.(type) {
    // 是一个响应收到的事件
        case *network.EventResponseReceived:
            resp := ev.Response
            if len(resp.Headers) != 0 {
                //将这个resp转成json
                response, _ := resp.MarshalJSON()
                var res = &UrlResponse{}
                json.Unmarshal(response, &res)
                // 我们只关心是图片地址的url
                if strings.Contains(res.Url, ".jpg") || strings.Contains(res.Url, "f=JPEG") {
                    // 去对每个图片地址下载图片
                    go download(res.Url)
            }
        }
    }
}

func download(url string) {
    tempFile, _ := ioutil.TempFile("", "*.jpg")
    resp := requests.Get(url, ext.Timeout(5))
    if resp == nil {
        return
    }
    resp.Save(tempFile.Name())
    fmt.Printf("已保存图片至 %s\n", tempFile.Name())
}
```

#### 使用 tools.FileDownload

```go
import (
	"github.com/sari3l/requests"
	"github.com/sari3l/requests/ext"
	"github.com/sari3l/requests/tools"
	"github.com/sari3l/requests/types"
)

func main() {
    headers := types.Dict{
        "User-Agent": tools.RandomUserAgent(),
    }
    resp := requests.Get("https://www.x.com/video?page=1", ext.Headers(headers))
    tools.FileDownload(resp, []string{"jpg"}, "")
}
```