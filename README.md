

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


## DB

```
docker compose exec postgres psql -U vegshare
```

```
docker compose down && docker compose up -d
docker compose logs -f <service_name>

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
