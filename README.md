А можешь переписать свой предыдущий ответ.Только также как и раньше погружай меня в контекст всего . 
в коде коментируй каждую строчку подробно чтобы я понимал откуда что берется, 
и как пекеты взаимодействуют между собой,  так как опыта у меня практически нет.

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
