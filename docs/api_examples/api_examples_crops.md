### Управление Культурами (Crops)

**1. Создание культуры (Требуются права администратора)**

```bash
ACCESS_TOKEN="your_admin_access_token"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Tomato", "description": "Red and juicy", "planting_time": 30, "harvest_time": 90}' http://localhost:8080/api/v1/catalog/crops
```

*Успешный ответ (201 Created):*
```json
{
    "id": "62d71460-4689-4e3d-8e17-101ade9ab271",
    "name": "Tomato",
    "description": "Red and juicy",
    "planting_time": 30,
    "harvest_time": 90,
    "created_at": "2025-08-20T20:27:26.441016618Z",
    "updated_at": "2025-08-20T20:27:26.441016618Z"
}
```

**2. Получение списка всех культур**

```bash
ACCESS_TOKEN="your_access_token"
curl -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/catalog/crops
```