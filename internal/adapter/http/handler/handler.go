package handler

import "net/http"

type SomeHandler struct {
}

func NewSomeHandler() *SomeHandler {
	return &SomeHandler{}
}

func (h *SomeHandler) Handle(w http.ResponseWriter, r *http.Request) {

}
