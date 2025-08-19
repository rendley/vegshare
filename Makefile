# Переменные
# Пути теперь указываются от корня проекта
DB_URL=postgres://vegshare:vegshare@localhost:5432/vegshare?sslmode=disable
MIGRATIONS_DIR=./backend/migrations
BIN_NAME=vegshare

.PHONY: all build run migrate-up migrate-down test clean

all: build

# Сборка
build:
	(cd backend && go build -o bin/$(BIN_NAME) ./cmd/main.go)

build-worker:
	(cd backend && go build -o bin/worker ./cmd/worker/main.go)

# Запуск
run:
	(cd backend && go run ./cmd/main.go)

run-worker:
	(cd backend && go run ./cmd/worker/main.go)

# Миграции
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-new:
	@read -p "Введите название миграции: " name; \
	migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) "$$name"

# Тесты
test:
	(cd backend && go test -v ./...)

# Очистка
clean:
	rm -rf backend/bin/
	go clean

# Установка зависимостей
deps:
	cd backend && go mod download
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Запуск с hot-reload (для разработки)
dev:
	(cd backend && air)

# Генерация документации
swag:
	swag init -g ./backend/cmd/main.go
