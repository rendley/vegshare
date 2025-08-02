VegShare/  
├── mobile/          # Мобильное приложение (React Native / Flutter / Kotlin + Swift)  
├
├── backend/             # Бекенд (микросервисы)  
│   ├── api-gateway/     # Единая точка входа (Nginx/Kong)  
│   ├── auth-service/    # Аутентификация (JWT/OAuth2)  
│   ├── user-service/    # Управление пользователями  
│   ├── farm-service/    # Логика фермы (урожай, животные)  
│   ├── payment-service/ # Платежи (Stripe/Yookassa)  
│   ├── notification-service/ # Уведомления (WebSocket, Firebase)  
│   └── ...  
├── infra/               # Инфраструктура  
│   ├── k8s/            # Kubernetes-манифесты  
│   ├── terraform/      # Облачная инфраструктура  
│   └── monitoring/     # Grafana, Prometheus  
└── docs/               # Документация (Swagger, архитектура)


## Auth-Service (Аутентификация)

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
│       ├── auth.go      # Обработчики auth (регистрация, логин)
│       └── types.go     # Структуры запросов/ответов
│   ├── server/         # Логика (регистрация, JWT)
│   │     └── server.go   # сервер
│   ├── service/         # Логика (регистрация, JWT)
│   ├── repository/      # Работа с БД
│   └── models/          # Сущности (User, Session)
├── pkg/
│   ├── config/          # Конфиги
│   ├    └── config.yaml # загрузчик конфигов
│   └── jwt/             # Генерация токенов
└── go.mod

```

