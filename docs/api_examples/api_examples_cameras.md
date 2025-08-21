### Управление Камерами (Cameras)

**1. Создание камеры для грядки**

Предполагается, что у вас уже есть ID грядки (`PLOT_ID`). 

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="510a4a6b-d279-46e8-8492-bce2bdf38197"
curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" \
-d '{"name": "Front View", "rtsp_path_name": "front_view_cam"}' \
http://localhost:8080/api/v1/farm/plots/$PLOT_ID/cameras
```

*Успешный ответ (201 Created):*
```json
{
    "id": "5ef864ef-fc68-4cef-81b0-115fa6ff3a19",
    "plot_id": "510a4a6b-d279-46e8-8492-bce2bdf38197",
    "name": "Front View",
    "rtsp_path_name": "front_view_cam",
    "created_at": "2025-08-21T21:46:48.624713934Z",
    "updated_at": "2025-08-21T21:46:48.624713934Z"
}
```

**2. Получение списка камер для грядки**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="510a4a6b-d279-46e8-8492-bce2bdf38197"
curl -s -X GET -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/plots/$PLOT_ID/cameras
```

*Успешный ответ (200 OK):*
```json
[
    {
        "id": "5ef864ef-fc68-4cef-81b0-115fa6ff3a19",
        "plot_id": "510a4a6b-d279-46e8-8492-bce2bdf38197",
        "name": "Front View",
        "rtsp_path_name": "front_view_cam",
        "created_at": "2025-08-21T21:46:48.624714Z",
        "updated_at": "2025-08-21T21:46:48.624714Z"
    }
]
```

**3. Удаление камеры**

```bash
ACCESS_TOKEN="your_access_token"
CAMERA_ID="5ef864ef-fc68-4cef-81b0-115fa6ff3a19"
curl -s -i -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/cameras/$CAMERA_ID
```

*Успешный ответ (204 No Content):*
```
HTTP/1.1 204 No Content
```

```