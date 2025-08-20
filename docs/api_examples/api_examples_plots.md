### Управление Грядками (Plots)

**1. Создание грядки в теплице**

```bash
ACCESS_TOKEN="your_access_token"
GREENHOUSE_ID="582b551b-f13f-4518-b398-3567f000a39a" # ID реальной теплицы
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Test Plot", "size": "1x1", "camera_url": "http://example.com/camera"}' http://localhost:8080/api/v1/farm/greenhouses/$GREENHOUSE_ID/plots
```

*Успешный ответ (201 Created):*
```json
{
    "id": "51175ae1-a6ae-45e2-9423-cce34fffcd63",
    "greenhouse_id": "582b551b-f13f-4518-b398-3567f000a39a",
    "name": "Test Plot",
    "size": "1x1",
    "status": "available",
    "camera_url": "http://example.com/camera",
    "created_at": "2025-08-20T20:24:11.658855201Z",
    "updated_at": "2025-08-20T20:24:11.658855201Z"
}
```

**2. Получение списка грядок для теплицы**

```bash
ACCESS_TOKEN="your_access_token"
GREENHOUSE_ID="582b551b-f13f-4518-b398-3567f000a39a" # ID реальной теплицы
curl -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/greenhouses/$GREENHOUSE_ID/plots
```

**3. Обновление грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="51175ae1-a6ae-45e2-9423-cce34fffcd63" # ID самой грядки
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Updated Plot", "size": "2x2", "status": "rented", "camera_url": "http://example.com/new_camera"}' http://localhost:8080/api/v1/farm/plots/$PLOT_ID
```

**4. Удаление грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="51175ae1-a6ae-45e2-9423-cce34fffcd63" # ID самой грядки
curl -i -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/plots/$PLOT_ID
```