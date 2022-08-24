# Tools

为方便使用，默认提供了一些小工具，目前分为四大类

- convert
- hash
- hooks
- ja3

## Convert

> func ConvertGbkToUtf8(str string) string

> func ConvertUtf8ToGbk(str string) string

> func CovertStructToJson(obj any) map[string]any

将结构体转换为`map[string]any`格式数据，需要注释有`json:"xxxx"`

> func ConvertStructToDict(obj any) ext.Dict

将结构体转换为`ext.Dict`格式数据，需要注释有`dict:"xxxx"`

> func ConvertValueToString(obj reflect.Value) string

获取反射值字符串内容

> func HookResponseGbkToUtf8(response any) (error, any)

转换response.Html编码为UTF-8，需自行hook至`requests.Response`

> func HookResponseUtf8ToGbk(response any) (error, any)

转换response.Html编码为GBK，需自行hook至`requests.Response`

## Hash

大多是比较常见的函数，具体不解释

> func HmacSha256(data []byte, secret []byte) []byte

> func HmacSha256Base64Encode(data []byte, secret []byte) string

> func Md5(data []byte) string

## Ja3

> func HookClientJA3Func(fingerprint string) types.Hook

接收指纹字符串，返回对应Hook函数，需自行hook至`http.client`

## CloudFlare

> func HookCloudFlareWorkFunc(workHost string, headers types.Dict) types.Hook

可用于绕过cloudflare，利用work反代进行访问，自动设置以下两个header

- Px-Host 会从原始请求中提取目标域名，可被显式覆盖
- Px-IP 自动生成IPv4，可被显示覆盖

最重要保证链接到 cloudflare 的出口是白的