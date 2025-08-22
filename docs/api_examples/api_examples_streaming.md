### Управление Стримингом (Streaming)

**1. Установка WebSocket-соединения для стрима**

Этот эндпоинт используется для установки WebSocket-соединения с сервером для последующей WebRTC-сигнализации. Для теста можно использовать `curl`, но для полноценной работы потребуется WebSocket-клиент.

```bash
ACCESS_TOKEN="your_access_token"
CAMERA_ID="85b7e63a-52c4-4225-8feb-9b2016bfb443"
curl -i -N \
  -H "Connection: Upgrade"
  -H "Upgrade: websocket"
  -H "Host: localhost:8080"
  -H "Origin: http://localhost:8080"
  -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ=="
  -H "Sec-WebSocket-Version: 13"
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  http://localhost:8080/api/v1/stream/$CAMERA_ID
```

*Успешный ответ (101 Switching Protocols):*
```
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```

