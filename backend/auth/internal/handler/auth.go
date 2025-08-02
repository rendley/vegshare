package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// homeHandler обрабатывает запросы к корневому пути ("/").
func (h *Handler) homeHandler(w http.ResponseWriter, r *http.Request) {
	//Проверяем что запрос именно GET
	if r.Method != http.MethodGet {
		// Если метод не GET, возвращаем ошибку 405 (Method Not Allowed).
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}
	// Пишем ответ в http.ResponseWriter.
	fmt.Fprintf(w, "Welcome to VegShare!")
}

// registerHandler обрабатывает POST запросы "/register"
func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Register request received")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed! Use POST method!")
		return
	}
	// fmt.Fprintf(w, "Register endpoint (will be implemented later)")
	var userID int
	err := h.db.QueryRow("INSERT INTO users (email) VALUES ($1) RETURNING id",
		"test@example.com").Scan(&userID)
	if err != nil {
		h.logger.Errorf("Database error: %v", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	h.logger.Infof("User regisered with ID: %d", userID)
	fmt.Fprintf(w, "User ID: %d", userID)

}

// loginHandler обрабатывает POST /login.
func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// fmt.Fprint(w, "Login endpoint (will be implemented later)")

	// 2. Читаем и парсим JSON-тело.
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}
	// 3. Валидация (минимальная).
	if req.Email == "" || len(req.Password) < 6 {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	// 4. Заглушка: проверка логина/пароля (позже заменим на БД).
	if req.Email != "test@example.com" || req.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 5. Формируем успешный ответ.
	response := LoginResponse{
		Token:  "fake-jwt-token",
		UserID: "123",
	}
	//6. Сериализуем ответ в JSON и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
