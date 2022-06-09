## 安装

通过go工具进行安装

```shell
go get github.com/sari3l/requests
```

## 简易使用

在项目中`import "github.com/sari3l/requests"`即可实现基础请求

```golang
import (
    "fmt"
    "github.com/sari3l/requests"
)

func main() {
    resp := requests.Get("https://github.com/")
    fmt.Println(resp.Content)
}
```

## 添加参数

需要额外`import "github.com/sari3l/requests/ext"`，可在[扩展参数](extensions.md?id=可选参数)中查看具体支持的参数内容

```golang
import (
    "fmt"
    "github.com/sari3l/requests"
    "github.com/sari3l/requests/ext"
)

func main() {
    cookies := ext.Dict{
        "key1": "value2",
    }
    
    resp := requests.Post("https://github.com/", ext.Cookies(cookies))
    fmt.Println(resp.Content)
}
```