### Получение списка всех культур

**Запрос:**

```bash
curl -i http://localhost:8080/api/v1/crops
```

**Успешный ответ (если культур нет):**

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 16 Aug 2025 17:52:00 GMT
Content-Length: 4

null
```
