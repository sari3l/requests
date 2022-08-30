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

### 常规方法

#### func Json()

返回`*gjson.Result`

#### func XPath()

返回`*parser.XpathNode`

#### func Text()

返回Document中的所有`Text`节点内容

#### func Save(path string)

此方法会将`Response.Raw`写入路径对应文件

#### func ContentType()

返回响应头中的`Content-Type`值

#### func URLs()

返回响应页面中的所有链接

### 动态渲染
#### func Render() *Response

动态渲染页面，将`<HTML>`内容写入`Response.HTML`属性

```go
resp := requests.Get("https://www.google.com", ext.Timeout(3)).Render()
fmt.Println(resp.Html)
```

#### func Snapshot(png bool) []byte

动态渲染页面后截图，返回`png`或`jpeg`对应`[]byte`数据，需要自行存储

```go
resp := requests.Get("https://www.google.com", ext.Timeout(3))
buf := resp.Snapshot(true)
if tmpFile, err := ioutil.TempFile("", "*.png"); err != nil {
    fmt.Println(err)
} else {
    tmpFile.Write(*buf)
    fmt.Println(tmpFile.Name())
}
```

#### func CustomRender(eventListener func(ctx context.Context), flags []chromedp.ExecAllocatorOption, actions ...chromedp.Action) *Response

自定义事件监听、无头参数设置，以及多组操作执行

```go
resp := requests.Get("https://www.google.com", ext.Timeout(3))
resp.CustomRender(nil, []chromedp.ExecAllocatorOption{chromedp.Flag("headless", false)}, chromedp.Sleep(1000*time.Second))
```