# Indonesia Life Online

Open World Multiplayer Life Simulation beridentitas Indonesia — dibangun dengan **Unreal Engine 5**, dedicated server, dan backend server-authoritative.

> Working title. Bukan GTA clone: tujuannya dunia virtual Indonesia yang **terasa hidup**.

---

## Status

| Item | Status |
|---|---|
| Fase | **0 → masuk 0.1 foundation** |
| MVP target | **0.1** (lihat docs) |
| Engine | Unreal Engine 5 |
| Backend | Go API scaffold (auth + character + migrations) |
| Infra lokal | Neon cloud Postgres (Docker optional) |
| Remote | https://github.com/AFR-projection/GTA |

---

## Visi singkat

Pemain login dan langsung merasa: **"Ini Indonesia banget."**

Hidup, kerja, bisnis, rumah, kendaraan, sosialisasi — cerita dibuat oleh pemain, bukan campaign linear. Detail suasana (warung, gang, sawah, tol, pasar, dll.) lebih penting daripada map raksasa kosong.

Dokumen lengkap: [`docs/VISION.md`](docs/VISION.md)

---

## Setup lokal (tanpa Docker)

Database pakai **Neon** (cloud). Docker lokal tidak wajib.

### 1) Install tools

- [Go](https://go.dev/dl/) 1.22+
- Akun [Neon](https://console.neon.tech) (Postgres gratis)
- Git
- Unreal Engine 5 (belakangan, setelah API hijau)

### 2) Environment

```powershell
cd c:\Users\User\Documents\aldopr
Copy-Item .env.example .env
```

Isi `DATABASE_URL` di `.env` dengan connection string Neon (`sslmode=require`).

### 3) Jalankan API

```powershell
cd backend
go mod tidy
go run .\cmd\api
```

Health check: `http://localhost:8080/healthz`

Detail: [`backend/README.md`](backend/README.md) · tugas harian: [`docs/ONBOARDING.md`](docs/ONBOARDING.md)

> `infra/docker-compose.yml` optional — hanya jika suatu saat mau Postgres/Redis lokal.

---

## Mulai dari sini (untuk kontributor)

Baca berurutan:

1. [`docs/PRD.md`](docs/PRD.md) — **Product Requirements** (apa & kenapa)  
2. [`docs/VISION.md`](docs/VISION.md) — visi kreatif jangka panjang  
3. [`docs/MVP.md`](docs/MVP.md) — **apa yang boleh dikerjakan di 0.1** (scope lock)  
4. [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — bagaimana sistem saling bicara  
5. [`docs/TECH_DECISIONS.md`](docs/TECH_DECISIONS.md) — stack yang dikunci & alasan  
6. [`docs/ROADMAP.md`](docs/ROADMAP.md) — urutan versi  
7. [`docs/ONBOARDING.md`](docs/ONBOARDING.md) — tugas konkret anak buah / kontributor baru  

**Aturan emas:** gameplay & correctness server-side dulu; grafik AAA final belakangan. Asset harus berlisensi / orisinal — jangan salin merek nyata.

---

## Arsitektur (ringkas)

```
UE5 Client  →  UE5 Dedicated Game Server  →  Go API  →  PostgreSQL (Neon)
                                              ↓
                                           Redis + R2 (storage)
```

Semua uang, inventory, dan ownership diputuskan di server/backend — client hanya kirim intent dan tampilkan hasil.

Detail: [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md)

---

## Struktur repo

```
aldopr/
├── README.md
├── .env.example
├── docs/
├── backend/           ← Go API (auth, characters, migrations)
├── game/              ← Unreal project (belum)
├── infra/             ← docker-compose.yml
└── tools/
```

---

## MVP 0.1 (ringkas)

Login, character creation, open world kecil, multiplayer, chat global, motor + mobil, inventory, uang, bank, SPBU, warung, rumah beli sederhana, save ke database, NPC dasar, siang/malam, cuaca. Target **20–50 pemain** / server.

Checklist penuh: [`docs/MVP.md`](docs/MVP.md)

---

## Prinsip tim

1. Gameplay > grafik  
2. Backend scalable sejak dini  
3. Sistem penting = server-side  
4. Optimasi jangan ditunda total ke akhir  
5. Asset bisa diganti tanpa rewrite gameplay  
6. Dunia hidup > dunia luas kosong  
7. Setiap update harus ada nilai buat pemain  
8. Hormati copyright — desain orisinal  
9. Kode modular & terdokumentasi  
10. Bangun untuk bertahun-tahun, bukan satu demo

---

## Kontribusi

- Conventional commits: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`
- Fitur di luar MVP 0.1 → diskusikan dulu; update `docs/MVP.md` / `docs/ROADMAP.md` jika disetujui
- Jangan commit secret (`.env`, key, token)

---

## Lisensi

TBD — tentukan sebelum distribusi publik / playtest eksternal.
