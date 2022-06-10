# Hook

Hook 用于在底层请求前后，对`请求`、`响应`内容进行修改，可以避免可选参数对某些设置的可能修改

在`ext.defaultHooksList`中定义了默认可 Hook 对象，分别为

- request -> requests.Request
- client -> http.Client
- response -> requests.Response

任意对象均可被多个Hook函数处理、修改

## Hook Func

所有`自定义Hook函数`均须满足`ext.Hook`函数类型

```go
type Hook func(object any) (error, any)
```

为了方便使用，如果是通过`反射`修改了 object，可以直接返回`nil, nil`

```go
return nil, nil
```

## Hook Dict

Hook字典是为方便使用创建的数据类型

- 键：`request/client/response`之一
- 值：`Hook Func`列表

```go
type HooksDict map[string][]Hook
```

## Hook Ext

自定义Hook函数通过`ext.Hooks`可选参数进行装填，单次请求进行Hook效率不高，推荐封装使用

```go
import (
    "fmt"
    "github.com/sari3l/requests"
    "github.com/sari3l/requests/ext"
)

func main() {
    hooks := ext.HooksDict{
        "response": []ext.Hook{printHeaders},
    }
    resp := requests.Get("https://www.google.com", ext.Hooks(hooks))
    fmt.Println(resp.Content)
}

func printHeaders(response any) (error, any) {
    fmt.Printf("%+v\n", response.(requests.Response).Header)
    return nil, response
}
```