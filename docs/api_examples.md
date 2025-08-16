# Примеры запросов к API VegShare

Этот файл содержит рабочие примеры `curl` запросов для тестирования и демонстрации API.

---

### Создание новой фермы

**Запрос:**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "Тестовая Ферма Gemini", "location": "Интернет"}' http://localhost:8080/api/v1/farms
```

**Успешный ответ (201 Created):**

```json
{
    "id": "14e9b6f2-c837-4d60-9e68-507e04ed0df6",
    "name": "Тестовая Ферма Gemini",
    "location": "Интернет"
}
```

---

### Получение списка всех ферм

**Запрос:**

```bash
curl -i http://localhost:8080/api/v1/farms
```

**Успешный ответ (если ферм нет):**

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 16 Aug 2025 16:51:30 GMT
Content-Length: 4

null
```

---

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
