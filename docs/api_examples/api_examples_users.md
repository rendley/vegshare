## User Registration and Profile

### Register a new user

**Request:**
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "name": "Test User",
  "email": "testuser@example.com",
  "password": "password123"
}' http://localhost:8080/api/v1/register
```

**Response:**
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user_id": "..."
}
```

### Get user profile

**Request:**
```bash
curl -H "Authorization: Bearer <YOUR_ACCESS_TOKEN>" http://localhost:8080/api/v1/users/me
```

**Response:**
```json
{
  "id": "...",
  "email": "testuser@example.com",
  "name": "Test User",
  "avatar_url": null,
  "created_at": "...",
  "updated_at": "..."
}
```
