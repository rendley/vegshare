### Управление Грядками (Plots)

**1. Создание грядки в теплице**

```bash
# Замените {greenhouseID} на ID реальной теплицы
curl -X POST -H "Content-Type: application/json" -d '{"name": "A-1", "size": "1x2m", "camera_url": "rtsp://1.2.3.4/stream1"}' http://localhost:8080/api/v1/greenhouses/41653306-f4a1-47ca-8418-1bd960433cc8/plots
```

*Успешный ответ (201 Created):*
```json
{
    "id": "872c5482-f284-43fc-88ab-607e93ba127f",
    "greenhouse_id": "41653306-f4a1-47ca-8418-1bd960433cc8",
    "name": "A-1",
    "size": "1x2m",
    "status": "available",
    "camera_url": "rtsp://1.2.3.4/stream1",
    "created_at": "2025-08-16T19:35:51.801992022Z",
    "updated_at": "2025-08-16T19:35:51.801992022Z"
}
```

**2. Получение списка грядок для теплицы**

```bash
# Замените {greenhouseID} на ID реальной теплицы
curl http://localhost:8080/api/v1/greenhouses/41653306-f4a1-47ca-8418-1bd960433cc8/plots
```

**3. Обновление грядки**

```bash
# Замените {plotID} на ID самой грядки
curl -X PUT -H "Content-Type: application/json" -d '{"name": "A-1", "size": "1x2m", "status": "maintenance", "camera_url": "rtsp://1.2.3.4/stream1"}' http://localhost:8080/api/v1/plots/872c5482-f284-43fc-88ab-607e93ba127f
```

**4. Удаление грядки**

```bash
# Замените {plotID} на ID самой грядки
curl -i -X DELETE http://localhost:8080/api/v1/plots/872c5482-f284-43fc-88ab-607e93ba127f
```
