### Управление Теплицами (Greenhouses)

**1. Создание теплицы на земельном участке**

```bash
# Замените {landParcelID} на ID реального участка
curl -X POST -H "Content-Type: application/json" -d '{"name": "Южная", "type": "Поликарбонат"}' http://localhost:8080/api/v1/land-parcels/5285869e-0c3e-4e44-abcc-8accf260affe/greenhouses
```

*Успешный ответ (201 Created):*
```json
{
    "id": "c7957095-e545-4926-93e1-682896e98247",
    "land_parcel_id": "5285869e-0c3e-4e44-abcc-8accf260affe",
    "name": "Южная",
    "type": "Поликарбонат",
    "created_at": "2025-08-16T19:02:42.132124523Z",
    "updated_at": "2025-08-16T19:02:42.132124523Z"
}
```

**2. Получение списка теплиц для участка**

```bash
# Замените {landParcelID} на ID реального участка
curl http://localhost:8080/api/v1/land-parcels/5285869e-0c3e-4e44-abcc-8accf260affe/greenhouses
```

**3. Обновление теплицы**

```bash
# Замените {id} на ID самой теплицы
curl -X PUT -H "Content-Type: application/json" -d '{"name": "Южная-1", "type": "Стекло"}' http://localhost:8080/api/v1/greenhouses/c7957095-e545-4926-93e1-682896e98247
```

**4. Удаление теплицы**

```bash
# Замените {id} на ID самой теплицы
curl -i -X DELETE http://localhost:8080/api/v1/greenhouses/c7957095-e545-4926-93e1-682896e98247
```