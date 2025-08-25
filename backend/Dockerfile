# --- Этап 1: Сборщик ---
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Копируем УЖЕ СКАЧАННУЮ утилиту migrate и делаем ее исполняемой
COPY migrate /usr/local/bin/migrate
RUN chmod +x /usr/local/bin/migrate

# Копируем файлы зависимостей и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной исходный код
COPY . .

# Собираем приложение
RUN go build -o /vegshare ./cmd/main.go


# --- Этап 2: Финальный образ ---
FROM alpine:latest

WORKDIR /app

# Копируем миграции из сборщика
COPY --from=builder /src/migrations ./migrations

# Копируем утилиту migrate из сборщика
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Копируем скомпилированное приложение из сборщика
COPY --from=builder /vegshare /vegshare

# Открываем порт
EXPOSE 8080

# Команда для запуска
CMD [ "/vegshare" ]