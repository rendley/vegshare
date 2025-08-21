### Управление Камерами (Cameras)

**1. Создание камеры для грядки**

Предполагается, что у вас уже есть ID грядки (`PLOT_ID`). 

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="83cd4784-4724-4756-a21c-33f2b0c74293"
curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" \
-d '{"name": "Refactor Front View", "rtsp_path_name": "refactor_front_cam"}' \
http://localhost:8080/api/v1/plots/$PLOT_ID/cameras
```

*Успешный ответ (201 Created):*
```json
{
    "id": "2febab9e-e149-40fd-ab71-6bb4c3640a7b",
    "plot_id": "83cd4784-4724-4756-a21c-33f2b0c74293",
    "name": "Refactor Front View",
    "rtsp_path_name": "refactor_front_cam",
    "created_at": "2025-08-21T23:25:09.339526047Z",
    "updated_at": "2025-08-21T23:25:09.339526047Z"
}
```

**2. Получение списка камер для грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="83cd4784-4724-4756-a21c-33f2b0c74293"
curl -s -X GET -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/plots/$PLOT_ID/cameras
```

*Успешный ответ (200 OK):*
```json
[
    {
        "id": "2febab9e-e149-40fd-ab71-6bb4c3640a7b",
        "plot_id": "83cd4784-4724-4756-a21c-33f2b0c74293",
        "name": "Refactor Front View",
        "rtsp_path_name": "refactor_front_cam",
        "created_at": "2025-08-21T23:25:09.339526Z",
        "updated_at": "2025-08-21T23:25:09.339526Z"
    }
]
```

**3. Удаление камеры**

```bash
ACCESS_TOKEN="your_access_token"
CAMERA_ID="2febab9e-e149-40fd-ab71-6bb4c3640a7b"
curl -s -i -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/cameras/$CAMERA_ID
```

*Успешный ответ (204 No Content):*
```
HTTP/1.1 204 No Content
```

