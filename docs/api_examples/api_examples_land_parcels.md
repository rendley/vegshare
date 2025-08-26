### Управление Земельными участками (Land Parcels)

**1. Создание участка в регионе (Требуются права администратора)**

```bash
ACCESS_TOKEN="your_admin_access_token"
REGION_ID="a0e726e5-209b-4446-86a9-bd199b26740f" # ID реального региона
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Test Land Parcel"}' http://localhost:8080/api/v1/farm/regions/$REGION_ID/land-parcels
```

*Успешный ответ (201 Created):*
```json
{
    "id": "26533b65-f6a2-4cd2-82a8-0b5ad9392cc5",
    "region_id": "a0e726e5-209b-4446-86a9-bd199b26740f",
    "name": "Test Land Parcel",
    "created_at": "2025-08-20T20:22:43.581013759Z",
    "updated_at": "2025-08-20T20:22:43.581013759Z"
}
```

**2. Получение списка участков для региона**

```bash
ACCESS_TOKEN="your_access_token"
REGION_ID="a0e726e5-209b-4446-86a9-bd199b26740f" # ID реального региона
curl -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/regions/$REGION_ID/land-parcels
```

**3. Обновление участка (Требуются права администратора)**

```bash
ACCESS_TOKEN="your_admin_access_token"
LAND_PARCEL_ID="26533b65-f6a2-4cd2-82a8-0b5ad9392cc5" # ID самого участка
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Updated Land Parcel"}' http://localhost:8080/api/v1/farm/land-parcels/$LAND_PARCEL_ID
```

**4. Удаление участка (Требуются права администратора)**

```bash
ACCESS_TOKEN="your_admin_access_token"
LAND_PARCEL_ID="26533b65-f6a2-4cd2-82a8-0b5ad9392cc5" # ID самого участка
curl -i -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/land-parcels/$LAND_PARCEL_ID
```
