# Backend API (Go)

Server-authoritative persistence layer for Indonesia Life Online.

## Endpoints (MVP slice)

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/healthz` | no | Health check |
| POST | `/v1/auth/register` | no | Create account + JWT |
| POST | `/v1/auth/login` | no | Login + JWT |
| GET | `/v1/me` | Bearer | Current account |
| GET | `/v1/characters` | Bearer | List characters |
| POST | `/v1/characters` | Bearer | Create character (max 1 in MVP) |
| GET | `/v1/characters/{id}` | Bearer | Get owned character |
| POST | `/v1/characters/{id}/bank/deposit` | Bearer | Cash → bank (audited) |
| POST | `/v1/characters/{id}/bank/withdraw` | Bearer | Bank → cash (audited) |

## Local run

Prerequisites: **Go 1.22+**, **Docker Desktop** (Postgres + Redis).

```powershell
# from repo root
Copy-Item .env.example .env
cd infra
docker compose up -d
cd ..\backend
go mod tidy
go run .\cmd\api
```

API default: `http://localhost:8080`

### Smoke test

```powershell
# register
curl -s -X POST http://localhost:8080/v1/auth/register -H "Content-Type: application/json" -d "{\"email\":\"dev@ilo.local\",\"password\":\"password123\",\"display_name\":\"Dev\"}"

# login (save token)
curl -s -X POST http://localhost:8080/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"dev@ilo.local\",\"password\":\"password123\"}"

# create character
curl -s -X POST http://localhost:8080/v1/characters -H "Authorization: Bearer YOUR_TOKEN" -H "Content-Type: application/json" -d "{\"name\":\"Budi\",\"gender\":\"male\",\"skin_tone\":3,\"hair_style\":1,\"face_preset\":0}"
```

## Layout

```
backend/
├── cmd/api/           # process entrypoint
├── internal/
│   ├── api/           # HTTP handlers
│   ├── auth/          # JWT + password hashing
│   ├── config/
│   ├── db/            # pool + goose migrations
│   ├── httpx/
│   ├── middleware/
│   ├── models/
│   └── store/         # SQL repository
└── migrations/        # versioned schema
```
