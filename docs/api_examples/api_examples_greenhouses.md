### Управление Теплицами (Greenhouses)

**1. Создание теплицы на земельном участке (Требуются права администратора)**

```bash
ACCESS_TOKEN="your_admin_access_token"
LAND_PARCEL_ID="26533b65-f6a2-4cd2-82a8-0b5ad9392cc5" # ID реального участка
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Test Greenhouse", "type": "glass"}' http://localhost:8080/api/v1/farm/land-parcels/$LAND_PARCEL_ID/greenhouses
```

*Успешный ответ (201 Created):*
```json
{
    "id": "582b551b-f13f-4518-b398-3567f000a39a",
    "land_parcel_id": "26533b65-f6a2-4cd2-82a8-0b5ad9392cc5",
    "name": "Test Greenhouse",
    "type": "glass",
    "created_at": "2025-08-20T20:23:25.812555644Z",
    "updated_at": "2025-08-20T20:23:25.812555644Z"
}
```

**2. Получение списка теплиц для участка**

```bash
ACCESS_TOKEN="your_access_token"
LAND_PARCEL_ID="26533b65-f6a2-4cd2-82a8-0b5ad9392cc5" # ID реального участка
curl -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/land-parcels/$LAND_PARCEL_ID/greenhouses
```

**3. Обновление теплицы (Требуются права администратора)**

```bash
ACCESS_TOKEN="your_admin_access_token"
GREENHOUSE_ID="582b551b-f13f-4518-b398-3567f000a39a" # ID самой теплицы
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"name": "Updated Greenhouse", "type": "polycarbonate"}' http://localhost:8080/api/v1/farm/greenhouses/$GREENHOUSE_ID
```

**4. Удаление теплицы (Требуются права администратора)**

```bash
ACCESS_TOKEN="your_admin_access_token"
GREENHOUSE_ID="582b551b-f13f-4518-b398-3567f000a39a" # ID самой теплицы
curl -i -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/farm/greenhouses/$GREENHOUSE_ID
```
