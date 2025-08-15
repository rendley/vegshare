// Пакет handler
package handler

import "net/http"

// RegisterRoutes регистрирует все эндпоинты для обработчика фермы.
// Он принимает в качестве аргумента mux (серверный мультиплексор), к которому привязывает пути.
func (h *FarmHandler) RegisterRouter(mux *http.ServeMux) {
	// Мы связываем путь "POST /farms" с нашим методом-обработчиком createFarm.
	// Теперь все POST-запросы на /farms будут обрабатываться этим методом.
	mux.HandleFunc("POST /api/v1/farms", h.CreateFarm)
}
