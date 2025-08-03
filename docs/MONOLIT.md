Примерная структура MVP монолита

backend/
├── cmd/
│   └── main.go                     # Инициализация всех компонентов
├── internal/
│   ├── auth/                       # Аутентификация
│   │   ├── handler/                # HTTP-хендлеры
│   │   │   ├── auth.go             # Роуты /register, /login
│   │   │   └── middleware.go       # JWT-мидлварь
│   │   ├── service/                # Бизнес-логика
│   │   │   └── auth_service.go     # Регистрация, валидация
│   │   └── repository/             # Работа с БД
│   │       └── auth_repository.go  # Запросы к auth_users
│   ├── user/                       # Профили
│   │   ├── handler/
│   │   │   └── profile.go          # Роуты /profile
│   │   └── service/
│   │       └── profile_service.go  # Логика профилей
│   ├── farm/                       # Ферма
│   │   ├── handler/
│   │   │   ├── crops.go            # Роуты /farm/crops
│   │   │   └── animals.go          # Роуты /farm/animals
│   │   └── service/
│   │       ├── crops_service.go    # Логика урожая
│   │       └── animals_service.go  # Логика животных
│   ├── streaming/                  # Стриминг
│   │   ├── handler/
│   │   │   └── stream.go           # Роуты /stream
│   │   └── service/
│   │       └── stream_service.go   # Логика стримов
│   └── api/                        # Общие компоненты
│       ├── router.go               # Объединение всех роутов
│       └── responses.go            # Форматы JSON-ответов
├── migrations/
│   ├── 0001_init_auth.up.sql       # auth_users
│   ├── 0002_init_farm.up.sql       # farm_crops
│   └── 0003_init_user.up.sql       # user_profiles
└── pkg/
├── database/
│   └── postgres.go             # Инициализация PG
├── config/
│   └── config.go               # Загрузка config.yaml
└── events/
├── producer.go             # Отправка событий
└── consumer.go             # Обработка событий