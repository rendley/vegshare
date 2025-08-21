### Управление Грядками (Plots)

**1. Создание грядки**

Создает новую грядку в указанной теплице.

```bash
ACCESS_TOKEN="your_access_token"
GREENHOUSE_ID="c8ddd5d7-234f-4481-9cc1-a2bdb6db36e7" # ID существующей теплицы
curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" \
-d '{"name": "Refactor Test Plot", "size": "3x3", "greenhouse_id": "'"$GREENHOUSE_ID"'"}' \
http://localhost:8080/api/v1/plots
```

*Успешный ответ (201 Created):*
```json
{
    "id": "83cd4784-4724-4756-a21c-33f2b0c74293",
    "greenhouse_id": "c8ddd5d7-234f-4481-9cc1-a2bdb6db36e7",
    "name": "Refactor Test Plot",
    "size": "3x3",
    "status": "available",
    "created_at": "2025-08-21T23:22:54.233982Z",
    "updated_at": "2025-08-21T23:22:54.233982Z"
}
```

**2. Получение списка грядок для теплицы**

```bash
ACCESS_TOKEN="your_access_token"
GREENHOUSE_ID="c8ddd5d7-234f-4481-9cc1-a2bdb6db36e7"
curl -s -X GET -H "Authorization: Bearer $ACCESS_TOKEN" "http://localhost:8080/api/v1/plots?greenhouse_id=$GREENHOUSE_ID"
```

*Успешный ответ (200 OK):*
```json
[
    {
        "id": "83cd4784-4724-4756-a21c-33f2b0c74293",
        "greenhouse_id": "c8ddd5d7-234f-4481-9cc1-a2bdb6db36e7",
        "name": "Refactor Test Plot",
        "size": "3x3",
        "status": "available",
        "created_at": "2025-08-21T23:22:54.233982Z",
        "updated_at": "2025-08-21T23:22:54.233982Z"
    }
]
```

**3. Обновление грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="83cd4784-4724-4756-a21c-33f2b0c74293"
curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" \
-d '{"name": "Updated Refactor Plot", "size": "4x4", "status": "maintenance"}' \
http://localhost:8080/api/v1/plots/$PLOT_ID
```

*Успешный ответ (200 OK):*
```json
{
    "id": "83cd4784-4724-4756-a21c-33f2b0c74293",
    "greenhouse_id": "c8ddd5d7-234f-4481-9cc1-a2bdb6db36e7",
    "name": "Updated Refactor Plot",
    "size": "4x4",
    "status": "maintenance",
    "created_at": "2025-08-21T23:22:54.233982Z",
    "updated_at": "2025-08-21T23:24:32.938061749Z"
}
```

**4. Удаление грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="83cd4784-4724-4756-a21c-33f2b0c74293"
curl -s -i -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/plots/$PLOT_ID
```

*Успешный ответ (204 No Content):*
```
HTTP/1.1 204 No Content
```