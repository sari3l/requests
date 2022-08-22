# Response

对http.Response进行封装，方便hook处理，同时添加一些属性、方法

```go
type Response struct {
    *http.Response
    Session *session
    cookies []*http.Cookie
    Ok      bool
    Raw     []byte
    Html    string
    History []*Response
    Time    int64
}
```

## 属性

### Ok 

返回判断请求是否成功且响应值是否为200

### Raw

返回完整Body内容

### HTML

返回完整HTML字符内容

### History

返回请求到最终响应的所有响应历史

### Time

返回请求到最终响应的总用时

## 方法

### Json()

返回`*gjson.Result`

### XPath()

返回`*parser.XpathNode`

### Text()

返回Document中的所有`Text`节点内容

### Save(path string)

此方法会将`Response.Raw`写入路径对应文件

### ContentType()

返回响应头中的`Content-Type`值

### URLs()

返回响应页面中的所有链接

### Render(useExtCookies bool)

动态渲染页面，将`<HTML>`内容写入`Response.HTML`属性

参数`useExtCookies`设置如下：
- `true` 使用扩展参数中的Cookies属性
- `false` 使用响应得到的cookie