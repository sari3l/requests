# Tools

为方便使用，默认提供了一些小工具

## CloudFlare

### func HookCloudFlareWorkerFunc(workerHost string, headers types.Dict) types.Hook

利用CloudFlare Workers反代进行请求，如若没有header身份校验，需自行hook至requests.Response

```go
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
    })},
}

resp := requests.Get("https://www.google.com", ext.Headers(headers), ext.Hooks(hooks))
fmt.Println(resp.Html)
```

## Convert

### func ConvertGbkToUtf8(str string) string

### func ConvertUtf8ToGbk(str string) string
 
### func CovertStructToJson(obj any) map[string]any

将结构体转换为`map[string]any`格式数据，需要注释有`json:"xxxx"`

```go
type testStruct struct {
    Name string `json:"name,omitempty"`
    Value string `json:"value,omitempty"`
}
test := testStruct{Name: "sari3l"}
json := tools.CovertStructToJson(test)
resp := requests.Post("http://httpbin.org/post", ext.Json(json), ext.Proxy("http://127.0.0.1:8080"))
fmt.Println(resp.Json())
```

### func ConvertStructToDict(obj any) ext.Dict

将结构体转换为`ext.Dict`格式数据，需要注释有`dict:"xxxx"`

```go
type testStruct struct {
    Name string `dict:"name,omitempty"`
    Value string `dict:"value,omitempty"`
}
test := testStruct{Name: "sari3l"}
params := tools.ConvertStructToDict(test)
resp := requests.Get("http://httpbin.org/get", ext.Params(params))
fmt.Println(resp.Json())
```

### func ConvertValueToString(obj reflect.Value) string

获取反射值字符串内容

### func HookResponseGbkToUtf8(response any) (error, any)

转换response.Html编码为UTF-8，需自行hook至`requests.Response`

### func HookResponseUtf8ToGbk(response any) (error, any)

转换response.Html编码为GBK，需自行hook至`requests.Response`

## Hash

大多是比较常见的函数

### func HmacSha256(data []byte, secret []byte) []byte

### func HmacSha256Base64Encode(data []byte, secret []byte) string

### func Md5(data []byte) string

## Ja3

### func HookClientJA3Func(fingerprint string) ext.Hook

接收指纹字符串，返回对应Hook函数，需自行hook至`http.client`

## Random

### func RandomIPv4() string

### func RandomIPv6() string

### func RandomUserAgent() string