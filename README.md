# 🧩 TaskAPI

A clean, testable REST API for managing tasks — built with **Go**, **chi**, and **PostgreSQL**.  
Designed to demonstrate professional backend architecture, layering, and testing patterns.

---

## 🚀 Overview

TaskAPI provides CRUD operations for tasks, with support for filtering, pagination, and multiple database adapters (raw SQL and GORM).  
It focuses on **clarity, modularity, and maintainability**, following clean architecture principles.

---

## 🗂️ Architecture

```graphql
cmd/
api/ # Server entrypoint (main.go)
internal/
api/ # HTTP routing, handlers, middleware (RequestID, Logger, Recoverer)
service/ # Business logic and application rules
repository/ # Interfaces + shared structs (ListFilter, Pagination)
postgres/ # SQL implementation using database/sql
gorm/ # Alternative GORM implementation
pkg/
models/ # Domain models (Task, TaskStatus, etc.)
migrations/ # SQL schema migrations
```


### Data Flow

```java
HTTP Request
↓
Router / Handlers
↓
Service Layer (TaskService)
↓
Repository Interface
↓
Postgres Adapter (Raw SQL or GORM)
↓
Database
```

---

## ⚙️ Tech Decisions & Trade-offs

| Choice | Reason |
|--------|---------|
| **chi router** | Lightweight, idiomatic, and composable middleware. |
| **Layered architecture** | Decouples logic from frameworks and persistence. |
| **Repository interface** | Allows switching between SQL and GORM implementations easily. |
| **Plain SQL migrations** | Explicit schema evolution with version control. |
| **Pagination with limit/offset** | Simple and reliable for moderate data sizes. |
| **Enums for status** | Prevents invalid task states at compile time. |

---

## 🧠 Core Concepts

- **Clean architecture** — separation of concerns between API, service, and data layers.  
- **Dependency inversion** — high-level code depends on interfaces, not implementations.  
- **Middleware pipeline** — includes request logging, recovery, and unique request IDs.  
- **Testing** — unit tests for service logic; optional integration tests for repositories.  
- **Dockerized environment** — PostgreSQL service managed through Docker Compose.  
- **Environment-driven config** — `.env` file loaded automatically at runtime.  

---

## 🌐 API Endpoints

Base URL: `http://localhost:{PORT}`

| Method | Endpoint | Description |
|--------|-----------|-------------|
| **GET** | `/healthz` | Health check endpoint. |
| **GET** | `/tasks` | List tasks (supports filters, search, pagination). |
| **GET** | `/tasks/{id}` | Retrieve a task by ID. |
| **POST** | `/tasks` | Create a new task. |
| **PUT** | `/tasks/{id}` | Update a task by ID. |
| **DELETE** | `/tasks/{id}` | Delete a task by ID. |

### Query Parameters for `/tasks`

| Name | Type | Description |
|------|------|-------------|
| `status` | string | Filter by status (`todo`, `in_progress`, `done`). |
| `search` | string | Search by keyword in title or description. |
| `limit` | int | Max results to return (default 20). |
| `offset` | int | Results offset for pagination (default 0). |

---

## 🧾 Example Requests & Responses

### Health
```http
GET /healthz
→ 200 OK
{
  "status": "ok"
}
```
Create Task
```http
POST /tasks
Content-Type: application/json

{
  "title": "Write integration tests",
  "description": "Add repository integration tests",
  "status": "todo"
}

→ 201 Created
{
  "id": "ed16a595-3b28-4365-b415-07db0e94f3f6",
  "title": "Write integration tests",
  "description": "Add repository integration tests",
  "status": "todo",
  "createdAt": "2025-10-24T15:03:00Z",
  "updatedAt": "2025-10-24T15:03:00Z"
}
```

List Tasks
```http
GET /tasks?status=in_progress&limit=5
→ 200 OK
[
  {
    "id": "c1a8…",
    "title": "Fix login bug",
    "description": "Handle user sessions correctly",
    "status": "in_progress",
    "createdAt": "2025-10-24T15:03:00Z",
    "updatedAt": "2025-10-24T17:28:16Z"
  }
]
```
Update Task
```http
PUT /tasks/{id}
Content-Type: application/json

{
  "status": "done"
}

→ 200 OK
{
  "id": "c1a8…",
  "title": "Fix login bug",
  "description": "Handle user sessions correctly",
  "status": "done",
  "updatedAt": "2025-10-24T17:40:00Z"
}
```
Delete Task
```http
DELETE /tasks/{id}
→ 204 No Content
```

# 🧪 Testing

Two categories of tests are implemented:

| Type	| Scope | Location |
|--------|-------|----------|
| Unit tests | Service logic | /internal/service |

The integration tests connect to a dedicated test database (e.g. taskapi_test) and automatically clean up data between runs.
(Integrations test to implement)
