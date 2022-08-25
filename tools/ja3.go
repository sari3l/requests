package tools

import (
	"github.com/CUCyber/ja3transport"
	"github.com/sari3l/requests/types"
	"net/http"
	"reflect"
)

func HookClientJA3Func(fingerprint string) types.Hook {
	return func(client any) (error, any) {
		c := client.(http.Client)
		tr, _ := ja3transport.NewTransport(fingerprint)
		reflect.ValueOf(&c).Elem().FieldByName("Transport").Set(reflect.ValueOf(tr))
		return nil, c
	}
}
