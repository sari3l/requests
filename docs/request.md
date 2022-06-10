# Request

实际只是对`http.Request`的单纯封装，方便`hook`处理，暂未有其他作用

```go
type Request struct {
    *http.Request
}
```