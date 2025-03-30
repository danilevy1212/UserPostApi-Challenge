# API Documentation

This API provides endpoints for managing users and posts.

---

## Health

### `GET /health`

**Success**:  
- `200 OK`
```json
{ "status": "OK" }
```

**Failure**:  
- `503 Service Unavailable`
```json
{ "status": "service unavailable" }
```

---

## Users

### `POST /users`

Creates a user.  
**Request**:
```json
{ "name": "John Doe", "email": "john@example.com" }
```

**Success**:
- `201 Created`
```json
{ "id": 1, "name": "John Doe", "email": "john@example.com" }
```

**Failure**:
- `409 Conflict`
```json
{ "error": "user already exists" }
```
- `422 Unprocessable Entity`
```json
{ "error": "bad entity" }
```
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

### `GET /users`

Fetch all users.  
**Success**:
- `200 OK`
```json
[ { "id": 1, "name": "John Doe", "email": "john@example.com" }, ... ]
```

**Failure**:
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

### `GET /users/{id}`

Fetch user by ID.

**Success**:
- `200 OK`
```json
{ "id": 1, "name": "John Doe", "email": "john@example.com" }
```

**Failure**:
- `400 Bad Request`
```json
{ "error": "invalid id" }
```
- `404 Not Found`
```json
{ "error": "user not found" }
```
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

### `PUT /users/{id}`

Update user by ID (full replacement).  
**Request**:
```json
{ "name": "New Name", "email": "new@example.com" }
```

**Success**:
- `200 OK`
```json
{ "id": 1, "name": "New Name", "email": "new@example.com" }
```

**Failure**:
- `400 Bad Request`
```json
{ "error": "invalid id" }
```
- `404 Not Found`
```json
{ "error": "user not found" }
```
- `409 Conflict`
```json
{ "error": "email already in use" }
```
- `422 Unprocessable Entity`
```json
{ "error": "bad entity" }
```
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

### `DELETE /users/{id}`

Delete user by ID.  
**Success**:
- `204 No Content`

**Failure**:
- `400 Bad Request`
```json
{ "error": "invalid id" }
```
- `404 Not Found`
```json
{ "error": "user not found" }
```
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

## Posts

### `POST /posts`

Creates a post.  
**Request**:
```json
{ "title": "Post Title", "content": "Some content", "user_id": 1 }
```

**Success**:
- `201 Created`
```json
{ "id": 1, "title": "Post Title", "content": "Some content", "user_id": 1 }
```

**Failure**:
- `409 Conflict`
```json
{ "error": "userID doesn't exist" }
```
- `422 Unprocessable Entity`
```json
{ "error": "bad entity" }
```
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

### `GET /posts`

Fetch all posts.  
**Success**:
- `200 OK`
```json
[ { "id": 1, "title": "...", "content": "...", "user_id": 1 }, ... ]
```

**Failure**:
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

### `GET /posts/{id}`

Fetch post by ID.  
**Success**:
- `200 OK`
```json
{ "id": 1, "title": "...", "content": "...", "user_id": 1 }
```

**Failure**:
- `400 Bad Request`
```json
{ "error": "invalid id" }
```
- `404 Not Found`
```json
{ "error": "post not found" }
```
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

### `PUT /posts/{id}`

Update post by ID.  
**Request**:
```json
{ "title": "Updated Title", "content": "Updated content", "user_id": 1 }
```

**Success**:
- `200 OK`
```json
{ "id": 1, "title": "Updated Title", "content": "Updated content", "user_id": 1 }
```

**Failure**:
- `400 Bad Request`
```json
{ "error": "invalid id" }
```
- `404 Not Found`
```json
{ "error": "post not found" }
```
- `409 Conflict`
```json
{ "error": "userID doesn't exist" }
```
- `422 Unprocessable Entity`
```json
{ "error": "bad entity" }
```
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

### `DELETE /posts/{id}`

Delete post by ID.  
**Success**:
- `204 No Content`

**Failure**:
- `400 Bad Request`
```json
{ "error": "invalid id" }
```
- `404 Not Found`
```json
{ "error": "post not found" }
```
- `503 Service Unavailable`
```json
{ "error": "service unavailable" }
```

---

## Assumptions & Limitations

- Email must be unique across all users.
- Once a post is created, its `user_id` is permanent (ownership does not change).
- Error feedback is minimal, not field-specific.
- Only full updates are supported (PUT).
- DB connection is assumed always necessary; otherwise returns `503`.
