### Управление Регионами (Regions)

**1. Создание региона**

```bash
ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYxMDg3MDQsImlhdCI6MTc1NjA3OTkwNCwic3ViIjoiY2YzMGU1NzAtMDhkNS00NjY1LWE3NjYtZTc4ZmQ2YWZmMGU5In0.nAZMy5d_rh-KhhsNvNRLIFOkuXSDrRLahvCE_SxIBUc"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Krasnodarskiy kray"}' http://localhost:8080/api/v1/farm/regions
```

*Успешный ответ (201 Created):*
```json
{
    "id": "a0e726e5-209b-4446-86a9-bd199b26740f",
    "name": "Krasnodarskiy kray",
    "created_at": "2025-08-20T20:21:09.466186678Z",
    "updated_at": "2025-08-20T20:21:09.466186678Z"
}
```

**2. Получение списка всех регионов**

```bash
ACCESS_TOKEN="your_access_token"
curl -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/regions
```

**3. Обновление региона**

```bash
ACCESS_TOKEN="your_access_token"
REGION_ID="a0e726e5-209b-4446-86a9-bd199b26740f"
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Tambovskaya oblast"}' http://localhost:8080/api/v1/farm/regions/$REGION_ID
```

**4. Удаление региона**

```bash
ACCESS_TOKEN="your_access_token"
REGION_ID="a0e726e5-209b-4446-86a9-bd199b26740f"
curl -i -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/regions/$REGION_ID
```