// Пакет handler содержит логику маршрутизации HTTP-запросов.
package handler

import (
	"fmt"
	"net/http"
)

// New создаёт и возвращает роутер (реализует http.Handler).
func New() http.Handler {
	// Создаем мультиплексор (роутер) из стандартной библиотеке
	mux := http.NewServeMux()

	// Регистрируем обработчики для путей:
	// - "/" → корневой путь
	// - "/register" → регистрация
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/register", registerHandler)

	return mux
}

// homeHandler обрабатывает запросы к корневому пути ("/").
func homeHandler(w http.ResponseWriter, r *http.Request) {
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
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed! Use POST method!")
		return
	}
	fmt.Fprintf(w, "Register endpoint (will be implemented later)")
}
