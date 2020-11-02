package domain

import (
	"net/http"
)

type IResource interface {
	Routes() *Routes
	Render(w http.ResponseWriter, req *http.Request, status int, v interface{})
}
