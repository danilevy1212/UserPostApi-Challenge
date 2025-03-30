# UserPost API

A RESTful API to manage users and their posts, built with Go using the Gin web framework and ent for ORM.

---

## 🚀 Getting Started

### Prerequisites

- Go >= 1.24
- PostgreSQL
- [direnv](https://direnv.net/) (only if using Nix with flakes enabled) or manual `.env` loading
- Docker (optional)

---

## 🛠 Setup

This project uses [ent](https://entgo.io/) for database modeling and code generation. The generated code lives under:
```
internal/database/repositories/postgresql/ent/
```
The schema definitions are located in:
```
internal/database/repositories/postgresql/ent/schema/
```

### Generate Ent Code

To generate code from your schema definitions:

```bash
go run -mod=mod entgo.io/ent/cmd/ent generate \
  --target=internal/database/repositories/postgresql/ent \
  ./internal/database/repositories/postgresql/ent/schema
```

To create a new schema:

```bash
go run -mod=mod entgo.io/ent/cmd/ent new --target internal/database/repositories/postgresql/ent/schema <SchemaName>
```

### 1. Clone the repo

```bash
git clone https://github.com/danilevy1212/UserPostApi-Challenge.git
cd UserPostApi-Challenge
```

### 2. Configure environment variables

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

If you're using Nix with flakes enabled, you can alternatively use `direnv` to load them:

```bash
direnv allow
```

Otherwise, use a tool like `dotenv`, or manually export the environment variables in your shell

## 🏃 Running the API

Before running the API, you should run the database migration:
```bash
go run ./cmd/migration
```

### Option 1: Run locally

```bash
go run ./cmd/api
```

### Option 2: Run with Docker Compose

```bash
docker compose up api
```

You can also bring up the PostgreSQL database separately with:

```bash
docker compose up database
```

---

## 🧪 Testing

The project includes both unit and integration tests.

To run all tests:

```bash
go test -v ./...
```

To update failing snapshot tests:
```bash
UPDATE_SNAPS=true go test ./...
```

Snapshot tests use [`go-snaps`](https://github.com/gkampitakis/go-snaps) and standard assertions use [`stretchr/testify`](https://github.com/stretchr/testify).

---

## 📄 API Documentation

Full API documentation with expected inputs, outputs, and error formats can be found at:

```
./docs/API.md
```

Postman and Bruno collections are also provided under:

```
./devtools/
```

---

## ☸️ Kubernetes Deployment

To run the API in a local Kubernetes cluster:

```bash
minikube start --driver=docker
```

Apply the configuration files:

```bash
kubectl apply -f k8s.yaml
kubectl apply -f k8s-secret.yaml # Use k8s-secret.example.yaml as template
```

Build the Docker image inside the minikube environment:

```bash
eval $(minikube docker-env)
docker build -f Dockerfile -t challenge-api:latest .
```

Forward the database port:

```bash
kubectl port-forward svc/db-service 5432:5432
```

Run DB migration:

```bash
go run ./cmd/migration
```

## ✅ Features

- Full CRUD for Users and Posts
- Cleanly separated layers (models, handlers, repository)
- Database migration via Go
- Snapshot testing for handlers (no full end-to-end integration tests by design)
- Structured logging with [zerolog](https://github.com/rs/zerolog)
  - Pretty output in development
  - JSON logs in production

---

## 📌 Notes

- Email is unique per user
- Once a post is created, its ownership (`user_id`) cannot be changed
- Minimal, non-field-specific error feedback

