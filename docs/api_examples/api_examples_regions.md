### Управление Регионами (Regions)

**1. Создание региона**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "Московская область"}' http://localhost:8080/api/v1/regions
```

*Успешный ответ (201 Created):*
```json
{
    "id": "502c2e4a-a96c-4beb-9ada-f25ca4b87dc5",
    "name": "Московская область",
    "created_at": "2025-08-16T19:02:25.05631734Z",
    "updated_at": "2025-08-16T19:02:25.05631734Z"
}
```

**2. Получение списка всех регионов**

```bash
curl http://localhost:8080/api/v1/regions
```

**3. Обновление региона**

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"name": "Тамбовская область"}' http://localhost:8080/api/v1/regions/502c2e4a-a96c-4beb-9ada-f25ca4b87dc5
```

**4. Удаление региона**

```bash
curl -i -X DELETE http://localhost:8080/api/v1/regions/502c2e4a-a96c-4beb-9ada-f25ca4b87dc5
```
