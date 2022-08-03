# Requests
[![](https://img.shields.io/github/license/sari3l/requests?style=flat-square)](https://github.com/sari3l/requests/blob/main/LICENSE)
[![](https://img.shields.io/badge/made%20by-sari3l-blue?style=flat-square)](https://github.com/sari3l)
[![](https://img.shields.io/github/go-mod/go-version/sari3l/requests?style=flat-square)](https://go.dev/)
[![](https://img.shields.io/github/v/tag/sari3l/requests?style=flat-square)](https://github.com/sari3l/requests)

[![Go Report Card](https://goreportcard.com/badge/github.com/sari3l/requests)](https://goreportcard.com/report/github.com/sari3l/requests)
[![CodeFactor](https://www.codefactor.io/repository/github/sari3l/requests/badge)](https://www.codefactor.io/repository/github/sari3l/requests)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsari3l%2Frequests.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsari3l%2Frequests?ref=badge_shield)

<h1 align="center"><img src="https://raw.githubusercontent.com/sari3l/requests/main/docs/static/logo.png" alt="Logo"/></h1>

## 安装

```shell
go get github.com/sari3l/requests
```

## Quick Demo

```golang
import (
    "fmt"
    "github.com/sari3l/requests"
    "github.com/sari3l/requests/ext"
    "github.com/sari3l/requests/types"
)

func main() {
    // Requests Bearer Token
    auth := ext.BasicAuth{Username: "o94KGT3MlbT...", Password: "fNbL2ukEGyvuGSM7bAuoq..."}
    data := types.Dict{
        "grant_type": "client_credentials",
    }
    resp := requests.Post("https://api.twitter.com/oauth2/token", ext.Auth(auth), ext.Data(data))
    
    // Requests with Twitter API 2.0
    if resp != nil && resp.Ok {
        fmt.Println(resp.Json())
        token := ext.BearerAuth{Token: resp.Json().Get("access_token").Str}
        resp2 := requests.Get("https://api.twitter.com/2/users/by/username/Sariel_D", ext.Auth(token))
        fmt.Println(resp2.Json())
    }
}
```

## 链接

- [说明文档](https://requests.sari3l.com)

## Licenses

[MIT License](https://github.com/sari3l/requests/blob/main/LICENSE)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsari3l%2Frequests.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsari3l%2Frequests?ref=badge_large)