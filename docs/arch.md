VegShare/  
├── mobile/          # Мобильное приложение (React Native / Flutter / Kotlin + Swift)  
├
├── backend/             # Бекенд (микросервисы)  
│   ├── api-gateway/     # Единая точка входа (Nginx/Kong)  
│   ├── auth-service/    # Аутентификация (JWT/OAuth2)  
│   ├── user-service/    # Управление пользователями  
│   ├── farm-service/    # Логика фермы (урожай, животные)
│   ├── streaming-service/ # Новый сервис стриминга (Go)
│   ├── payment-service/ # Платежи (Stripe/Yookassa) на будущее
│   ├── notification-service/ # Уведомления (WebSocket, Firebase)  
│   └── ...  
├── infra/               # Инфраструктура  
│   ├── k8s/            # Kubernetes-манифесты  
│   ├── terraform/      # Облачная инфраструктура  
│   └── monitoring/     # Grafana, Prometheus  
└── docs/               # Документация (Swagger, архитектура)


## 1. Auth Service (чистая аутентификация)

#### вход/регистрация.

① Структура сервиса и как его пишем 
1.Настройка конфига → 2. Сервер → 3. Роутеры → 4. Обработчики

1. Планирование обработчика /login

Перед кодом определим требования:
    Метод: POST (логин — это изменение состояния).
    Данные: email и password в теле запроса (JSON).
    Валидация: Проверка формата email и длины пароля.
    Ответ: Успех (200) или ошибка (400, 401).


```
auth-service/
├── cmd/
│   └── main.go          # Точка входа
├── configs/
│   └── config.yaml      # конфигурация YAML
├── internal/
│   ├── handler/         # HTTP-обработчики и роутеры
│       ├── handler.go   # Регистрация роутов
│       └── helpers.go   # Вспомогательные методы
│       ├── auth.go      # Обработчики auth (регистрация, логин)
│       └── types.go     # Структуры запросов/ответов
│   ├── server/         # Логика (регистрация, JWT)
│   │     └── server.go   # сервер
│   ├── service/         # Логика (регистрация, JWT)
│   │   └── auth_service.go # Бизнес-логика
│   ├── repository/      # Работа с БД
│       └── auth_repo.go    # Работа с БД
│   └── models/          # Сущности (User, Session)
├── pkg/
│   ├── config/          # Конфиги
│   ├    └── config.yaml # загрузчик конфигов
│   └── jwt/             # Генерация токенов
│   └── security/
│        ├── hasher.go         # Интерфейс PasswordHasher
│        ├── bcrypt_hasher.go  # Реализация BcryptHasher
│        └── argon2_hasher.go  # Альтернативная реализация
└── go.mod

```

## 2. User-service → Хранение профилей.


## 3. Farm-service → Основная логика приложения.


## 4.Архитектура Streaming-Service
   Технологии:

Язык: Go (высокая производительность для стриминга).

Протокол: WebSocket (для реального времени) + HTTP/2 (для gRPC).

Брокер сообщений: NATS (если нужна масштабируемость).

Хранение видео: S3/MinIO + HLS/DASH для адаптивного стриминга.

Структура сервиса:

```
streaming-service/  
├── cmd/                  # Точка входа (main.go)  
├── internal/  
│   ├── delivery/         # Транспортные слои  
│   │   ├── http/         # REST/WebSocket handlers  
│   │   └── grpc/         # gRPC сервер (если нужно)  
│   ├── service/          # Бизнес-логика  
│   │   ├── stream.go     # Управление стримами  
│   │   └── transcoder.go # Транскодирование видео  
│   └── repository/       # Работа с данными  
│       ├── s3.go         # Запись в S3  
│       └── redis.go      # Кеш состояний стримов  
├── pkg/                  # Общие утилиты  
│   ├── config/           # Конфиги  
│   └── logger/           # Логирование  
├── proto/                # gRPC-контракты (stream.proto)  
├── Dockerfile            # Контейнеризация  
└── go.mod                # Зависимости  
```