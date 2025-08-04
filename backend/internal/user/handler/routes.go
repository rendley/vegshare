package handler

import "net/http"

func (h *UserHandler) RegisterRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/me", h.GetUserHandler)
	mux.HandleFunc("PATCH /users/me", h.UpdateUserHandler)
	mux.HandleFunc("DELETE /users/me", h.DeleteUserHandler)
}

//func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
//	// Применяем auth middleware ко всем user-роутам
//	mux.HandleFunc("GET /users/me", h.withAuth(h.GetUserHandler))
//	mux.HandleFunc("PATCH /users/me", h.withAuth(h.UpdateUserHandler))
//	mux.HandleFunc("DELETE /users/me", h.withAuth(h.DeleteUserHandler))
//}
//
//func (h *UserHandler) withAuth(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		if !h.auth.IsAuthenticated(r) {
//			h.respondWithError(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//		next(w, r)
//	}
//}

// Модерационные ендпонты  добавить

//func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
//	// Персональные эндпоинты
//	mux.HandleFunc("GET /users/me", h.withAuth(h.GetCurrentUser))
//	mux.HandleFunc("PATCH /users/me", h.withAuth(h.UpdateProfile))
//	mux.HandleFunc("DELETE /users/me", h.withAuth(h.DeactivateAccount))
//
//	// Модерационные эндпоинты
//	mux.HandleFunc("GET /users/{id}", h.withAdmin(h.GetUserByID))
//	mux.HandleFunc("PATCH /users/{id}", h.withAdmin(h.UpdateUserByID))
//	mux.HandleFunc("DELETE /users/{id}", h.withAdmin(h.DeactivateUserByID))
//
//	// Админ-листинг
//	mux.HandleFunc("GET /users", h.withAdmin(h.ListUsers))
//}
