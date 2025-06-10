# Go Todo List API

A secure, containerized REST API for managing a Todo list, built with Go, Chi, PostgreSQL, and documented using Swagger.

---

## ✨ Features

- Full CRUD for tasks (Create, Read, Update, Delete)
- Per-user task isolation via JWT authentication
- PostgreSQL backend with auto-migration
- Swagger/OpenAPI documentation at `/swagger/index.html`
- Health check and root endpoint
- Clean middleware chain: security headers, request logging, JWT auth
- Docker and Docker Compose support for easy setup

---

## 📦 Project Structure

```

├── Dockerfile                  # Multi-stage build
├── docker-compose.yaml         # Compose setup for API + PostgreSQL
├── docs/                       # Swagger doc files
├── handlers/                   # Route handler functions
├── middleware/                 # Auth, security, and logging middleware
├── models/                     # DB models and persistence logic
├── utils/                      # Helper utilities (JWT, etc.)
├── main.go                     # App entry point
├── go.mod / go.sum             # Go modules
└── README.md                   # You're here!

```

---

## 🚀 Getting Started

### ✅ Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/)
- Go (optional, only if running without Docker)

---

### ⚙️ Environment Variables

These environment variables configure the database connection:

| Variable      | Description                | Example    |
| ------------- | -------------------------- | ---------- |
| `DB_USER`     | PostgreSQL username        | `postgres` |
| `DB_PASSWORD` | PostgreSQL password        | `postgres` |
| `DB_NAME`     | Database name              | `tododb`   |
| `DB_HOST`     | Hostname of DB container   | `db`       |
| `DB_PORT`     | Port PostgreSQL listens on | `5432`     |

Defined in `docker-compose.yaml` and used internally by the app. You can override these variables in your local environment or `.env` file if needed.

---

### 🐳 Run with Docker Compose

1. Clone the repo:

```bash
git clone https://github.com/your-username/go-todo-list.git
cd go-todo-list
````

2. Build and start the services:

```bash
docker-compose up --build
```

3. Access the API:

   * Base URL: `http://localhost:8080`
   * Swagger UI: `http://localhost:8080/swagger/index.html`

---

### 💻 Run Locally Without Docker

> Make sure a PostgreSQL instance is running and reachable.

1. Export the necessary environment variables (or create a `.env` file and load it):

```bash
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=tododb
export DB_HOST=localhost
export DB_PORT=5432
```

2. Start the server:

```bash
go run main.go
```

---

## 🔐 Authentication

All `/tasks` endpoints require a **valid JWT token**.

* The API provides user registration and login endpoints to create accounts and obtain JWT tokens.

* To authenticate, send the token in the `Authorization` header as:

  ```
  Authorization: Bearer <your-token>
  ```

* `userID` is extracted from the JWT and used to isolate tasks per user.

* Public endpoints (no auth required):

  * `/`
  * `/health`
  * `/swagger/*`
  * `/users/register`
  * `/users/login`

---

## 📘 API Endpoints

| Method | Endpoint          | Description          | Auth Required |
| ------ | ----------------- | -------------------- | ------------- |
| GET    | `/`               | Welcome message      | ❌             |
| GET    | `/health`         | Health check         | ❌             |
| GET    | `/swagger/*`      | Swagger UI/docs      | ❌             |
| POST   | `/users/register` | User registration    | ❌             |
| POST   | `/users/login`    | User login (get JWT) | ❌             |
| GET    | `/tasks`          | List all tasks       | ✅             |
| POST   | `/tasks`          | Create a new task    | ✅             |
| GET    | `/tasks/{id}`     | Get task by ID       | ✅             |
| PUT    | `/tasks/{id}`     | Update task by ID    | ✅             |
| DELETE | `/tasks/{id}`     | Delete task by ID    | ✅             |

---

## 📄 Swagger / OpenAPI

Auto-generated Swagger docs are served at:

```
http://localhost:8080/swagger/index.html
```

Use it to test and explore your API.

---

## 🧪 Notes

* Database tables auto-migrate on startup.
* Tasks are isolated by user ID from JWT — each user only sees their own tasks.
* Health check and root endpoints are unauthenticated.
* Security headers are added globally via middleware.
* Graceful shutdown is handled on `SIGINT` / `SIGTERM`.

---

## 📚 License

This project is open-source and available under the [MIT License](./LICENSE.md).

