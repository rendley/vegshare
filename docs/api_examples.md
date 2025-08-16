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
    "id": "10d2b717-fb0b-40a3-9d44-3fe10b8ad746",
    "name": "Воронежская область",
    "created_at": "2025-08-16T18:15:55.084151016Z",
    "updated_at": "2025-08-16T18:15:55.084151016Z"
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
        "id": "10d2b717-fb0b-40a3-9d44-3fe10b8ad746",
        "name": "Воронежская область",
        "created_at": "2025-08-16T18:15:55.084151Z",
        "updated_at": "2025-08-16T18:15:55.084151Z"
    }
]
```

**3. Обновление региона**

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"name": "Тамбовская область"}' http://localhost:8080/api/v1/regions/10d2b717-fb0b-40a3-9d44-3fe10b8ad746
```

*Успешный ответ (200 OK):*
```json
{
    "id": "10d2b717-fb0b-40a3-9d44-3fe10b8ad746",
    "name": "Тамбовская область",
    "created_at": "2025-08-16T18:15:55.084151Z",
    "updated_at": "2025-08-16T18:23:35.216298Z"
}
```

**4. Удаление региона**

```bash
curl -i -X DELETE http://localhost:8080/api/v1/regions/10d2b717-fb0b-40a3-9d44-3fe10b8ad746
```

*Успешный ответ:*
```http
HTTP/1.1 204 No Content
```

---

### Управление Земельными участками (Land Parcels)

**1. Создание участка в регионе**

```bash
# Замените {regionID} на ID реального региона
curl -X POST -H "Content-Type: application/json" -d '{"name": "Черноземье"}' http://localhost:8080/api/v1/regions/10d2b717-fb0b-40a3-9d44-3fe10b8ad746/land-parcels
```

*Успешный ответ (201 Created):*
```json
{
    "id": "374c9b23-7c5b-4c36-b270-c6be37e6d3f3",
    "region_id": "10d2b717-fb0b-40a3-9d44-3fe10b8ad746",
    "name": "Черноземье",
    "created_at": "2025-08-16T18:23:14.858805413Z",
    "updated_at": "2025-08-16T18:23:14.858805413Z"
}
```

**2. Получение списка участков для региона**

```bash
# Замените {regionID} на ID реального региона
curl http://localhost:8080/api/v1/regions/10d2b717-fb0b-40a3-9d44-3fe10b8ad746/land-parcels
```

*Успешный ответ:*
```json
[
    {
        "id": "374c9b23-7c5b-4c36-b270-c6be37e6d3f3",
        "region_id": "10d2b717-fb0b-40a3-9d44-3fe10b8ad746",
        "name": "Черноземье",
        "created_at": "2025-08-16T18:23:14.858805Z",
        "updated_at": "2025-08-16T18:23:14.858805Z"
    }
]
```

**3. Обновление участка**

```bash
# Замените {id} на ID самого участка
curl -X PUT -H "Content-Type: application/json" -d '{"name": "Солнечная долина"}' http://localhost:8080/api/v1/land-parcels/374c9b23-7c5b-4c36-b270-c6be37e6d3f3
```

*Успешный ответ (200 OK):*
```json
{
    "id": "374c9b23-7c5b-4c36-b270-c6be37e6d3f3",
    "region_id": "10d2b717-fb0b-40a3-9d44-3fe10b8ad746",
    "name": "Солнечная долина",
    "created_at": "2025-08-16T18:23:14.858805Z",
    "updated_at": "2025-08-16T18:23:35.216298275Z"
}
```

**4. Удаление участка**

```bash
# Замените {id} на ID самого участка
curl -i -X DELETE http://localhost:8080/api/v1/land-parcels/374c9b23-7c5b-4c36-b270-c6be37e6d3f3
```

*Успешный ответ:*
```http
HTTP/1.1 204 No Content
```