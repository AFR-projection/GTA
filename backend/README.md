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
| GET | `/v1/shops/warung` | Bearer | Katalog warung |
| GET | `/v1/characters/{id}/inventory` | Bearer | List inventory |
| POST | `/v1/characters/{id}/shops/warung/purchase` | Bearer | Beli item (cash + inventory + audit) |
| GET | `/v1/housing/listings` | Bearer | Daftar rumah dijual |
| GET | `/v1/characters/{id}/houses` | Bearer | Rumah milik karakter |
| POST | `/v1/characters/{id}/houses/buy` | Bearer | Beli rumah (1x/karakter MVP) |

## Local run (recommended: Neon, no Docker)

Prerequisites: **Go 1.22+**, akun **Neon** (Postgres).

1. Buat project di https://console.neon.tech → copy connection string  
2. Root repo: `Copy-Item .env.example .env` → paste URI ke `DATABASE_URL`  
3. Jalankan:

```powershell
cd backend
go mod tidy
go run .\cmd\api
```

API default: `http://localhost:8080`

Migrasi jalan otomatis saat API start.

### Optional: Docker lokal

Kalau suatu saat mau DB di laptop: `cd infra && docker compose up -d` lalu pakai URL localhost di `.env`. **Tidak wajib.**

### Smoke test (PowerShell)

Di Windows PowerShell, jangan pakai `curl -H ...` (itu alias lain). Pakai:

```powershell
# register
Invoke-RestMethod -Method POST -Uri http://localhost:8080/v1/auth/register -ContentType 'application/json' -Body '{"email":"dev@ilo.local","password":"password123","display_name":"Dev"}'

# login → simpan token
$login = Invoke-RestMethod -Method POST -Uri http://localhost:8080/v1/auth/login -ContentType 'application/json' -Body '{"email":"dev@ilo.local","password":"password123"}'
$token = $login.token

# create character
Invoke-RestMethod -Method POST -Uri http://localhost:8080/v1/characters -Headers @{ Authorization = "Bearer $token" } -ContentType 'application/json' -Body '{"name":"Budi","gender":"male","skin_tone":3,"hair_style":1,"face_preset":0}'

# beli di warung (ganti CHARACTER_ID)
$charId = "CHARACTER_ID"
Invoke-RestMethod -Method POST -Uri "http://localhost:8080/v1/characters/$charId/shops/warung/purchase" -Headers @{ Authorization = "Bearer $token" } -ContentType 'application/json' -Body '{"item_key":"kopi_tubruk","quantity":2}'
```

Atau pakai `curl.exe` (bukan `curl`) kalau mau syntax curl klasik.

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
