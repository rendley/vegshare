package handler

import "net/http"

func (h *Handler) RegisterRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.homeHandler)
	mux.HandleFunc("POST /register", h.registerHandler)
	mux.HandleFunc("POST /login", h.loginHandler)
}
