package types

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

type AuthInter interface {
	Format(p any) error
}

type BasicAuth struct {
	Username string
	Password string
}

func (a BasicAuth) Format(p any) error {
	identity := a.Username + ":" + a.Password
	p.(*http.Header).Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(identity))))
	return nil
}

type BearerAuth struct {
	Token string
}

func (a BearerAuth) Format(p any) error {
	p.(*http.Header).Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	return nil
}
