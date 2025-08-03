# Переменные
DB_URL=postgres://vegshare:vegshare@localhost:5432/vegshare?sslmode=disable
MIGRATIONS_DIR=./migrations
BIN_NAME=vegshare

.PHONY: all build run migrate-up migrate-down test clean

all: build

# Сборка
build:
	go build -o bin/$(BIN_NAME) ./cmd/main.go

# Запуск
run:
	go run ./cmd/main.go

# Миграции
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-new:
	@read -p "Введите название миграции: " name; \
	migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) "$$name"

# Тесты
test:
	go test -v ./...

# Очистка
clean:
	rm -rf bin/
	go clean

# Установка зависимостей
deps:
	go mod download
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Запуск с hot-reload (для разработки)
dev:
	air

# Генерация документации
swag:
	swag init -g cmd/main.go