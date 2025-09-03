### Управление Cтроениями (Structures)

Это новые эндпоинты после рефакторинга `greenhouses` в `structures`.

**Переменные окружения для тестов:**
```bash
# Установите ваш токен администратора
ADMIN_TOKEN="your_admin_access_token"

# ID земельного участка из предустановленных данных
LAND_PARCEL_ID="26533b65-f6a2-4cd2-82a8-0b5ad9392cc5" 
```

---

**1. Создание строения на земельном участке (Требуются права администратора)**

*Запрос:*
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" \
-d '{"name": "Главная теплица", "type": "glass"}' \
http://localhost:8080/api/v1/farm/land-parcels/$LAND_PARCEL_ID/structures
```

*Пример успешного ответа (201 Created):*
```json
{
    "id": "e2a6b3c8-a1b2-4c3d-8e4f-a1b2c3d4e5f6",
    "land_parcel_id": "26533b65-f6a2-4cd2-82a8-0b5ad9392cc5",
    "name": "Главная теплица",
    "type": "glass",
    "created_at": "2025-09-03T10:00:00Z",
    "updated_at": "2025-09-03T10:00:00Z"
}
```
**Сохраните `id` созданной структуры для следующих шагов.**
```bash
STRUCTURE_ID="e2a6b3c8-a1b2-4c3d-8e4f-a1b2c3d4e5f6"
```

---

**2. Получение списка строений для участка**

*Запрос:*
```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
http://localhost:8080/api/v1/farm/land-parcels/$LAND_PARCEL_ID/structures
```

---

**3. Обновление строения (Требуются права администратора)**

*Запрос:*
```bash
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" \
-d '{"name": "Обновленная Главная теплица", "type": "polycarbonate"}' \
http://localhost:8080/api/v1/farm/structures/$STRUCTURE_ID
```

---

**4. Создание грядки в строении (Требуются права администратора)**

*Запрос:*
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" \
-d '{"structure_id": "'$STRUCTURE_ID'", "name": "Грядка А-1", "description": "Под томаты", "size": {"width": 1.2, "length": 5.0}}' \
http://localhost:8080/api/v1/plots
```

*Пример успешного ответа (201 Created):*
```json
{
    "id": "f7b1a2c3-d4e5-f6a7-b8c9-d0e1f2a3b4c5",
    "structure_id": "e2a6b3c8-a1b2-4c3d-8e4f-a1b2c3d4e5f6",
    "name": "Грядка А-1",
    "description": "Под томаты",
    "size": {
        "width": 1.2,
        "length": 5
    },
    "created_at": "2025-09-03T10:05:00Z",
    "updated_at": "2025-09-03T10:05:00Z"
}
```
**Сохраните `id` созданной грядки.**
```bash
PLOT_ID="f7b1a2c3-d4e5-f6a7-b8c9-d0e1f2a3b4c5"
```

---

**5. Удаление грядки (Требуются права администратора)**

*Запрос:*
```bash
curl -i -X DELETE -H "Authorization: Bearer $ADMIN_TOKEN" \
http://localhost:8080/api/v1/plots/$PLOT_ID
```
*Ожидаемый ответ: `HTTP/1.1 204 No Content`*

---

**6. Удаление строения (Требуются права администратора)**

*Запрос:*
```bash
curl -i -X DELETE -H "Authorization: Bearer $ADMIN_TOKEN" \
http://localhost:8080/api/v1/farm/structures/$STRUCTURE_ID
```
*Ожидаемый ответ: `HTTP/1.1 204 No Content`*

---

**7. Проверка удаления строения**

*Запрос:*
```bash
curl -i -H "Authorization: Bearer $ADMIN_TOKEN" \
http://localhost:8080/api/v1/farm/structures/$STRUCTURE_ID
```
*Ожидаемый ответ: `HTTP/1.1 404 Not Found`*
