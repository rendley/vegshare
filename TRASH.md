
След шаги

Реализовать JWT-токены
Добавить миграции для БД
Написать интеграционные тесты
Добавить валидацию email
Реализовать rate limiting

Добавить валидацию (например, github.com/go-playground/validator).
Подключить БД (PostgreSQL):
Создать UserRepository.
Реализовать проверку пароля.
Генерация JWT (например, github.com/golang-jwt/jwt).
Middleware для логирования и аутентификации.

Добавить зависимости (БД, Redis) в Server.
Реализовать бизнес-логику в internal/service.
Подключить middleware (логгирование, CORS).

Добавим стандартный роутер.
Реализуем обработчики (/register, /login).

Подключить БД (сохранять пользователей).

import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Строка-ключ: все запросы, начинающиеся с /api, будут проксироваться
      '/api': {
        // Цель: наш бэкенд-сервер
        target: 'http://localhost:8080',
        // Изменяем origin, чтобы бэкенд думал, что запрос пришел с того же хоста
        changeOrigin: true,
      },
    }
  }
})

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Mount("/auth", s.AuthHandler.Routes())
		r.Mount("/catalog", s.CatalogHandler.Routes())

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(s.mw.AuthMiddleware)
			r.Mount("/users", s.UserHandler.Routes())
			r.Mount("/farm", s.FarmHandler.Routes())
			r.Mount("/leasing", s.LeasingHandler.Routes())
			r.Mount("/operations", s.OperationsHandler.Routes())
		})
	})