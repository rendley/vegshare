### Управление Земельными участками (Land Parcels)

**1. Создание участка в регионе**

```bash
# Замените {regionID} на ID реального региона
curl -X POST -H "Content-Type: application/json" -d '{"name": "Солнечногорский"}' http://localhost:8080/api/v1/regions/502c2e4a-a96c-4beb-9ada-f25ca4b87dc5/land-parcels
```

*Успешный ответ (201 Created):*
```json
{
    "id": "5285869e-0c3e-4e44-abcc-8accf260affe",
    "region_id": "502c2e4a-a96c-4beb-9ada-f25ca4b87dc5",
    "name": "Солнечногорский",
    "created_at": "2025-08-16T19:02:35.690003717Z",
    "updated_at": "2025-08-16T19:02:35.690003717Z"
}
```

**2. Получение списка участков для региона**

```bash
# Замените {regionID} на ID реального региона
curl http://localhost:8080/api/v1/regions/502c2e4a-a96c-4beb-9ada-f25ca4b87dc5/land-parcels
```

**3. Обновление участка**

```bash
# Замените {id} на ID самого участка
curl -X PUT -H "Content-Type: application/json" -d '{"name": "Солнечная долина"}' http://localhost:8080/api/v1/land-parcels/5285869e-0c3e-4e44-abcc-8accf260affe
```

**4. Удаление участка**

```bash
# Замените {id} на ID самого участка
curl -i -X DELETE http://localhost:8080/api/v1/land-parcels/5285869e-0c3e-4e44-abcc-8accf260affe
```