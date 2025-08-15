# Запускаем базу

```
docker compose up -d
```

# Установите migrate и air (для hot-reload)

```
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/cosmtrek/air@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#  Проводим миграции

```
make migrate-down && make migrate-up
```

# Запускаем проект

 ```
make run
 ```


## DB

```
docker compose exec postgres psql -U vegshare

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

Проверка создания

SELECT column_name, data_type FROM information_schema.columns 
WHERE table_name = 'users';

docker compose exec postgres bash



docker compose exec postgres psql -U vegshare -c "CREATE TABLE users (id SERIAL PRIMARY KEY, email TEXT UNIQUE);"

```

```
docker compose down && docker compose up -d
docker compose logs -f <service_name>


docker compose down -v && docker compose up -d

docker compose down -v  # Удаляет контейнеры и volumes
docker compose up -d    # Запускает с новой версией Postgres 16 и пустой БД


# Для тестового окружения (удалить данные и начать заново):
docker compose down -v && docker compose up -d

# Для продакшена (миграция через дамп):
docker compose exec postgres pg_dumpall -U vegshare > dump.sql
docker compose down -v
# Обновите docker-compose.yml до postgres:16
docker compose up -d
cat dump.sql | docker compose exec -T postgres psql -U vegshare
```
