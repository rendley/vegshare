## Таблица для auth

Make

# Установите migrate и air (для hot-reload)
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/cosmtrek/air@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

```
make migrate-down && make migrate-up
```

backend/
└── migrations/
├── 000001_init_auth.up.sql       # Создаёт users, refresh_tokens
├── 000001_init_auth.down.sql     # Удаляет их
├── 000002_add_profiles.up.sql    # Добавляет profiles (позже)
└── 000002_add_profiles.down.sql  # Удаляет profiles

## Стратегия развития БД
Сейчас:

Заливаем только auth-таблицы (users, refresh_tokens).

Добавляем модуль user:

Создаём 000002_add_profiles.up.sql с полями profiles.

Добавляем модуль farm:

000003_add_farms.up.sql — таблица farms + связи.

Если нужно пересоздать БД:

```
make migrate-down && make migrate-up
```

```

-- Таблица users (основа для аутентификации)
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица refresh_tokens (для JWT)
CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индекс для ускорения поиска пользователей
CREATE INDEX idx_users_email ON users(email);
```

users and farms

```
-- Таблица profiles (будет расширяться в модуле user)
CREATE TABLE profiles (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    full_name   VARCHAR(100),
    avatar_url  TEXT
);

-- Таблица farms (заготовка для модуля farm)
CREATE TABLE farms (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

```