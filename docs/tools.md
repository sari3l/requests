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

## Hash

大多是比较常见的函数，具体不解释

> func HmacSha256(data string, secret string) []byte

> func HmacSha256Base64Encode(data string, secret string) string

> func Md5(str string) string

## Hooks

> func HookResponseGbkToUtf8(response any) (error, any)

转换response.Content编码为UTF-8，需自行hook至`requests.Response`

> func HookResponseUtf8ToGbk(response any) (error, any)

转换response.Content编码为GBK，需自行hook至`requests.Response`

## Ja3

> func HookClientJA3Func(fingerprint string) ext.Hook

接收指纹字符串，返回对应Hook函数，需自行hook至`http.client`
