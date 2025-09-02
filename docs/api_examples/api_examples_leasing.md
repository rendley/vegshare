### Управление Арендой (Leasing)

**1. Аренда юнита (грядки, теплицы и т.д.)**

```bash
ACCESS_TOKEN="your_access_token"
UNIT_ID="51175ae1-a6ae-45e2-9423-cce34fffcd63" # ID объекта (например, грядки)

curl -X POST -H "Authorization: Bearer $ACCESS_TOKEN" -H "Content-Type: application/json" \
-d '{
  "unit_id": "'"$UNIT_ID"'",
  "unit_type": "plot"
}' \
http://localhost:8080/api/v1/leasing
```

*Успешный ответ (201 Created):*
```json
{
    "id": "4bbd645c-675c-49db-ab0e-618f4054ba77",
    "unit_id": "51175ae1-a6ae-45e2-9423-cce34fffcd63",
    "unit_type": "plot",
    "user_id": "a0ee6e8e-747c-4442-b111-abced36a3824",
    "start_date": "2025-08-20T20:25:34.896932Z",
    "end_date": "2025-11-20T20:25:34.896932Z",
    "status": "active",
    "created_at": "2025-08-20T20:25:34.896932Z",
    "updated_at": "2025-08-20T20:25:34.896932Z"
}
```

**2. Получение списка своих аренд**

```bash
ACCESS_TOKEN="your_access_token"
curl -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/leasing
```

*Успешный ответ (200 OK):*
```json
[
    {
        "id": "4bbd645c-675c-49db-ab0e-618f4054ba77",
        "unit_id": "51175ae1-a6ae-45e2-9423-cce34fffcd63",
        "unit_type": "plot",
        "user_id": "a0ee6e8e-747c-4442-b111-abced36a3824",
        "start_date": "2025-08-20T20:25:34.896932Z",
        "end_date": "2025-11-20T20:25:34.896932Z",
        "status": "active",
        "created_at": "2025-08-20T20:25:34.896932Z",
        "updated_at": "2025-08-20T20:25:34.896932Z"
    }
]
```