# REST API

A small Todo REST API built with Go as a practice project for learning REST APIs, HTTP handlers, SQLite, and basic project structure.

## Features

- CRUD operations for todos
- SQLite with [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) — pure Go, no CGO
- Request validation with [`go-playground/validator`](https://github.com/go-playground/validator)
- YAML configuration with [`cleanenv`](https://github.com/ilyakaznacheev/cleanenv)
- Consistent JSON responses
- Graceful shutdown with a 5-second timeout

## Requirements

- Go 1.27.1+

## Getting Started

```bash
git clone https://github.com/MHNahib/go-rest-api-basic.git
cd go-rest-api-basic
go mod download
mkdir storage
go run ./cmd/rest-api
```

The server runs on `localhost:8080` by default.

## Build

```bash
go build -o bin/rest-api ./cmd/rest-api
./bin/rest-api -config ./config/local.config.yaml
```

## Configuration

Config path priority:

1. `-config` flag
2. `CONFIG_PATH` environment variable
3. `./config/local.config.yaml`

Example:

```yaml
env: "dev"

server:
  host: "localhost"
  port: 8080
  address: "localhost:8080"

database:
  type: "sqlite"
  path: "storage/db.sqlite"
```

The `todos` table is created automatically on startup.

## API

| Method | Endpoint       | Description  |
| ------ | -------------- | ------------ |
| GET    | `/`            | Root         |
| GET    | `/health`      | Health check |
| GET    | `/todo`        | List todos   |
| POST   | `/todo`        | Create todo  |
| PUT    | `/todo`        | Update todo  |
| POST   | `/todo/delete` | Delete todo  |

Example todo:

```json
{
  "id": 1,
  "text": "Clean the room",
  "done": false
}
```

`text` is required for create/update. `id` is required for update/delete.

### Examples

```bash
# Create
curl -X POST http://localhost:8080/todo \
  -H "Content-Type: application/json" \
  -d '{"text":"Clean the room","done":false}'

# List
curl http://localhost:8080/todo

# Update
curl -X PUT http://localhost:8080/todo \
  -H "Content-Type: application/json" \
  -d '{"id":1,"text":"Clean the room","done":true}'

# Delete
curl -X POST http://localhost:8080/todo/delete \
  -H "Content-Type: application/json" \
  -d '{"id":1}'
```

See [`api.http`](api.http) for ready-to-use requests.

## Project Structure

```text
.
├── cmd/rest-api/          # Application entry point
├── config/                # Configuration
├── internal/
│   ├── config/            # Config loading
│   ├── http/handlers/     # HTTP handlers
│   ├── storage/           # Storage interface & SQLite
│   ├── types/             # Models
│   └── utils/response/    # JSON responses
├── api.http               # Example requests
└── storage/               # SQLite database
```

Handlers depend on the `storage.Storage` interface, keeping the HTTP layer independent from the database implementation.
