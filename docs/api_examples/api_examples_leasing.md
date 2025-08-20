### Управление Арендой (Leasing)

**1. Аренда грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="51175ae1-a6ae-45e2-9423-cce34fffcd63" # ID грядки
curl -X POST -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/leasing/plots/$PLOT_ID/lease
```

*Успешный ответ (201 Created):*
```json
{
    "id": "4bbd645c-675c-49db-ab0e-618f4054ba77",
    "plot_id": "51175ae1-a6ae-45e2-9423-cce34fffcd63",
    "user_id": "a0ee6e8e-747c-4442-b111-abced36a3824",
    "start_date": "2025-08-20T20:25:34.896932189Z",
    "end_date": "2025-11-20T20:25:34.896932189Z",
    "status": "active",
    "created_at": "2025-08-20T20:25:34.896932189Z",
    "updated_at": "2025-08-20T20:25:34.896932189Z"
}
```

**2. Получение списка своих аренд**

```bash
ACCESS_TOKEN="your_access_token"
curl -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/leasing/me/leases
```
