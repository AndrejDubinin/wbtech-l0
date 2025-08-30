package http

import "net/http"

type IndexHandler struct{}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	GetSuccessResponseWithBody(w, []byte("Service L0: Order is online"))
}
