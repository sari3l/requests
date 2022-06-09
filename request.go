package requests

import (
	"net/http"
)

type Request struct {
	*http.Request
}
