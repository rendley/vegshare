# Примеры запросов к API VegShare

Этот файл содержит рабочие примеры `curl` запросов для тестирования и демонстрации API.

---

### Получение списка всех культур

**Запрос:**

```bash
curl -i http://localhost:8080/api/v1/crops
```

**Успешный ответ (если культур нет):**

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 16 Aug 2025 17:52:00 GMT
Content-Length: 4

null
```

---

### Управление Регионами (Regions)

**1. Создание региона**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "Воронежская область"}' http://localhost:8080/api/v1/regions
```

*Успешный ответ (201 Created):*
```json
{
    "id": "bb9dc28c-b189-4d64-8b7b-4d62773dac3f",
    "name": "Воронежская область",
    "created_at": "2025-08-16T18:01:32.851958402Z",
    "updated_at": "2025-08-16T18:01:32.851958402Z"
}
```

**2. Получение списка всех регионов**

```bash
curl http://localhost:8080/api/v1/regions
```

*Успешный ответ:*
```json
[
    {
        "id": "bb9dc28c-b189-4d64-8b7b-4d62773dac3f",
        "name": "Воронежская область",
        "created_at": "2025-08-16T18:01:32.851958Z",
        "updated_at": "2025-08-16T18:01:32.851958Z"
    }
]
```

**3. Обновление региона**

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"name": "Тамбовская область"}' http://localhost:8080/api/v1/regions/bb9dc28c-b189-4d64-8b7b-4d62773dac3f
```

*Успешный ответ (200 OK):*
```json
{
    "id": "bb9dc28c-b189-4d64-8b7b-4d62773dac3f",
    "name": "Тамбовская область",
    "created_at": "2025-08-16T18:01:32.851958Z",
    "updated_at": "2025-08-16T18:02:08.326344267Z"
}
```

**4. Удаление региона**

```bash
curl -i -X DELETE http://localhost:8080/api/v1/regions/bb9dc28c-b189-4d64-8b7b-4d62773dac3f
```

*Успешный ответ:*
```http
HTTP/1.1 204 No Content
Date: Sat, 16 Aug 2025 18:02:15 GMT
```