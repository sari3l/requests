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
    session := requests.Session(5, "", true, true)
    _, prep := requests.PrepareRequest("get", "https://ja3er.com/json", nil, nil, nil, nil, nil, nil, nil, nil, nil)
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