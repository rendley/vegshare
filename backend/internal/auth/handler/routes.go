package handler

import "net/http"

func (h *AuthHandler) RegisterRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.homeHandler)
	mux.HandleFunc("POST /register", h.registerHandler)
	mux.HandleFunc("POST /login", h.loginHandler)
}
