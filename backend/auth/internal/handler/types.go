package handler

// LoginRequest — структура для парсинга JSON-тела запроса.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse — структура для успешного ответа.
type LoginResponse struct {
	Token  string `json:"token"`   // JWT-токен (пока заглушка)
	UserID string `json:"user_id"` // ID пользователя
}
