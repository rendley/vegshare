### Управление Грядками (Plots)

**1. Создание грядки в теплице**

```bash
ACCESS_TOKEN="your_access_token"
GREENHOUSE_ID="9bd1f3fa-1e9e-4548-8fc6-df1c96be6723" # ID реальной теплицы
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Test Plot 1", "size": "2x2"}' http://localhost:8080/api/v1/farm/greenhouses/$GREENHOUSE_ID/plots
```

*Успешный ответ (201 Created):*
```json
{
    "id": "f6a11800-86d6-4b73-afa4-69d6e0553fab",
    "greenhouse_id": "9bd1f3fa-1e9e-4548-8fc6-df1c96be6723",
    "name": "Test Plot 1",
    "size": "2x2",
    "status": "available",
    "created_at": "2025-08-21T20:39:38.247821Z",
    "updated_at": "2025-08-21T20:39:38.247821Z"
}
```

**2. Получение списка грядок для теплицы**

```bash
ACCESS_TOKEN="your_access_token"
GREENHOUSE_ID="9bd1f3fa-1e9e-4548-8fc6-df1c96be6723" # ID реальной теплицы
curl -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/greenhouses/$GREENHOUSE_ID/plots
```

*Успешный ответ (200 OK):*
```json
[
    {
        "id": "f6a11800-86d6-4b73-afa4-69d6e0553fab",
        "greenhouse_id": "9bd1f3fa-1e9e-4548-8fc6-df1c96be6723",
        "name": "Updated Test Plot",
        "size": "2x3",
        "status": "maintenance",
        "created_at": "2025-08-21T20:39:38.247821Z",
        "updated_at": "2025-08-21T20:40:46.525396Z"
    }
]
```

**3. Обновление грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="f6a11800-86d6-4b73-afa4-69d6e0553fab" # ID самой грядки
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Updated Test Plot", "size": "2x3", "status": "maintenance"}' http://localhost:8080/api/v1/farm/plots/$PLOT_ID
```

*Успешный ответ (200 OK):*
```json
{
    "id": "f6a11800-86d6-4b73-afa4-69d6e0553fab",
    "greenhouse_id": "9bd1f3fa-1e9e-4548-8fc6-df1c96be6723",
    "name": "Updated Test Plot",
    "size": "2x3",
    "status": "maintenance",
    "created_at": "2025-08-21T20:39:38.247821Z",
    "updated_at": "2025-08-21T20:40:46.525396127Z"
}
```

**4. Удаление грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="f6a11800-86d6-4b73-afa4-69d6e0553fab" # ID самой грядки
curl -i -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/plots/$PLOT_ID
```

*Успешный ответ (204 No Content):*
```
HTTP/1.1 204 No Content
```
