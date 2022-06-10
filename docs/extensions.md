# 扩展参数

由于golang不支持可选参数，所以此项目通过抽象参数为Func进行处理，下面仍以可选参数来介绍，但请勿混淆概念

```go
type Ext func(ep *extensions?Package)
```

目前支持的可选参数有

- [ext.AllowRedirects](extensions?id=extallowredirectsbool)
- [ext.Auth](extensions?id=extauthextauthinter)
    - ext.BasicAuth
    - ext.BearerAuth
- [ext.Cookies](extensions?id=extcookiesextdict)
- [ext.Data](extensions?id=extdataextdict)
- [ext.Files](extensions?id=extfilesextdict)
- [ext.Headers](extensions?id=extheadersextdict)
- [ext.Hooks](extensions?id=exthooksexthookdict)
- [ext.Json](extensions?id=extjsonmapstringinterface)
- [ext.Params](extensions?id=extparamsextdict)
- [ext.Proxy](extensions?id=extproxystring)
- [ext.Stream](extensions?id=extstreamioreader)
- [ext.Verify](extensions?id=extverifybool)

为了使用可选参数，需要在文件中

```go
import "github.com/sari3l/requests/ext"
```

另外为了方便处理数据，对以下数据类型取了别名，可通过`ext.Dict`、`ext.List`调用

```go
type Dict map[string]string
type List []string
```

注：单个请求可设置多个可选参数，下面是对单个参数的解释，所以均只设置相关参数

## ext.AllowRedirects(bool)

> 是否开启自动跳转

默认`true`，即自动处理跳转至最终页面，同时会将中间响应保存在`Response.History`中

```go
var resp *requests.Response

resp = requests.Get("https://httpbin.org/redirect/2", ext.AllowRedirects(false))
fmt.Println(resp.StatusCode)
resp = requests.Get("https://httpbin.org/redirect/2", ext.AllowRedirects(true))
fmt.Println(resp.StatusCode)
```

## ext.Auth(ext.AuthInter)

Auth认证稍微有些特别，因为其多样性，所以其是以接口形式定义，具体实现为

- ext.BasicAuth
- ext.BearerAuth

```go
var auth ext.AuthInter
var resp *requests.Response

auth = ext.BasicAuth{Username: "test", Password: "test"}
resp = requests.Get("https://github.com", ext.Auth(auth))
fmt.Println(resp.Content)

auth = ext.BearerAuth{Token: "test"}
resp = requests.Get("https://httpbin.org/bearer", ext.Auth(auth))
fmt.Println(resp.Json())
```

## ext.Cookies(ext.Dict)

实际的Cookies并不友好，所以这里采用`ext.Dict`方便设置，在内部自动转换为`[]*http.Cookie`

```go
cookies := ext.Dict{
    "key": "value",
}
resp := requests.Get("https://httpbin.org/cookies", ext.Cookies(cookies))
fmt.Println(resp.Json())
```

## ext.Data(ext.Dict)

data内容最终转换为`*io.ReadCloser`数据，并会自动设置`Content-Type`为`application/x-www-form-urlencoded`

注：不会判断请求方法是否合理，需要自行注意

```go
data := ext.Dict{
    "key": "value",
}
resp := requests.Post("https://httpbin.org/post", ext.Data(data))
fmt.Println(resp.Json())
```

## ext.Files(ext.Dict)

files内容最终转换为`*io.ReadCloser`数据，并会自动设置`Content-Type`为`multipart/form-data`

- 键：文件名
- 值：文件所在绝对路径

```go
files := ext.Dict{
    "xxx.jpg": "/path/xxx.jpg",
}
resp := requests.Post("https://httpbin.org/post", ext.Files(files))
fmt.Println(resp.Json())
```

## ext.Headers(ext.Dict)

headers内容最终转化为`*http.Header`数据，在设置前会检查是否有非法值

```go
headers := ext.Dict{
    "Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
    "Accept-Encoding": "gzip, deflate, br",
    "Accept-Language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
    "Cache-Control":   "no-cache",
    "Connection":      "keep-alive",
    "User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36",
}
resp := requests.Get("https://httpbin.org/headers", ext.Headers(headers))
fmt.Println(resp.Json())
```

## ext.Hooks(ext.HookDict)

Hook相关内容稍微复杂，具体内容请看`指南`[Hook](hook.md)一节，这里只简单演示如何使用

```go
func main() {
    hooks := ext.HooksDict{
        "response": []ext.Hook{printHeaders},
    }
    resp := requests.Get("https://httpbin.org/headers", ext.Hooks(hooks))
    fmt.Println(resp.Json()
}

func printHeaders(response any) (error, any) {
    fmt.Printf("%+v\n", response.(requests.Response).Header)
    return nil, response
}
```

## ext.Json(map[string]interface{})

json实在没有直白一点的实现，所以目前采用`map[string]interface{}`，最终转换为`*io.ReadCloser`数据，并会自动设置`Content-Type`为`application/json`

```go
json := map[string]interface{}{
    "string": "test",
    "list":   []interface{}{"1", 2},
    "dict": map[string]interface{}{
        "key": "value",
    },
}
resp := requests.Post("https://httpbin.org/post", ext.Json(json))
fmt.Println(resp.Json())
```

## ext.Params(ext.Dict)

与直接在URL中拼接参数不同，通过`ext.Params`填充的参数会经过`URLEncode`

```go
var resp *requests.Response

params := ext.Dict{
    "key": "%%25",
}
resp = requests.Get("https://httpbin.org/get", ext.Params(params))
fmt.Println(resp.Json())

resp = requests.Get("https://httpbin.org/get?key=%%25")
fmt.Println(resp.Json())
```

## ext.Proxy(string)

与python-requests中proxy不同，`net.http.Transport.Proxy`只支持单条`url.URL`，所以需要自行确认代理协议

```go
resp := requests.Get("https://github.com/", ext.Proxy("http://127.0.0.1:8080"))
```

## ext.Stream(io.Reader)

> 数据需要采用流的形式进行传输

会自动设置`Content-Type`为`application/octet-stream`

```go
stream, _ := os.Open("/path/s.png")
resp := requests.Post("https://httpbin.org/post", ext.Stream(stream))
fmt.Println(resp.Json())
```

## ext.Timeout(int)

> 连接超时时限

```go
resp := requests.Get("https://httpbin.org/get", ext.Timeout(3))
fmt.Println(resp.Json())
```


## ext.Verify(bool)

> 是否校验证书合法性

默认`true`，当设置代理抓包时证书可能无法通过此校验需要设置为`false`

```go
resp := requests.Get("https://github.com/", ext.Verify(false))
```
