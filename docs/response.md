# Response

对http.Response进行封装，方便hook处理，同时添加一些属性、方法

```go
type Response struct {
    *http.Response
    cookies []*http.Cookie
    Ok      bool
    Raw     []byte
    Html    string
    History []*Response
    Time    int64
}
```

## Ok 

返回判断请求是否成功且响应值是否为200

## Raw

返回完整Body内容

## HTML

返回完整HTML字符内容

## History

返回请求到最终响应的所有响应历史

## Time

返回请求到最终响应的总用时

## Json()

返回`*gjson.Result`

## XPath()

返回`*parser.XpathNode`

## Text()

返回Document中的所有`Text`节点内容

## Save(string)

此方法会将`Response.Raw`写入路径对应文件

## ContentType()

返回响应头中的`Content-Type`值

## URLs()

返回响应页面中的所有链接