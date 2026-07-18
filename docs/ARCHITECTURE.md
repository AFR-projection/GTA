# Architecture

Arsitektur target untuk Indonesia Life Online. Prinsip utama: **server-authoritative**, **scalable foundation**, **modular**.

```
┌─────────────────┐     ┌──────────────────────┐     ┌─────────────────┐
│  UE5 Client     │────▶│  Dedicated Game      │────▶│  Backend API    │
│  (Windows)      │◀────│  Server (UE5)        │◀────│  (Go)           │
└─────────────────┘     └──────────────────────┘     └────────┬────────┘
                                                              │
                                    ┌─────────────────────────┼─────────────────────────┐
                                    ▼                         ▼                         ▼
                             PostgreSQL                  Redis                   Object Storage
                             (Neon)                      (session/              (R2 / setara)
                                                         presence)               avatar, uploads
```

---

## Lapisan sistem

### 1. Game Client (Unreal Engine 5)
- Rendering, input, UI, audio, prediksi gerakan lokal
- Mengirim **intent** (mau beli, mau drive, mau chat) — bukan hasil final
- Menerima state yang di-authorize server
- Tidak menyimpan truth untuk uang, inventory, ownership

### 2. Dedicated Game Server (UE5)
- Authority untuk: posisi relevan, vehicle possession, interaksi world, combat/physics yang perlu fairness
- Replikasi ke client
- Memanggil Backend API untuk operasi persistent (uang, inventory, rumah, kendaraan ownership)
- Target kapasitas awal: 20–50 pemain / instance

### 3. Backend API (Go)
- REST untuk auth, character, catalog, admin/tools
- WebSocket opsional untuk event non-gameplay (notifikasi, presence hub)
- JWT authentication
- Validasi seluruh aksi ekonomi & inventory
- Rate limiting, audit log transaksi

### 4. PostgreSQL (Neon)
Source of truth persistensi:

- Account, Character
- Inventory, Economy (wallet/bank), Transaction history
- Vehicle, House, Business (nanti)
- Quest, Achievement, Guild, Friends, Marketplace (nanti)
- Progress metadata

### 5. Redis
- Session / online player set
- Cooldown & rate-limit counters
- Leaderboard sementara / realtime ephemeral state
- Matchmaking queue (fase berikutnya)

### 6. Object Storage (Cloudflare R2 atau setara)
- Avatar, screenshot, replay, file upload
- DB hanya menyimpan metadata / URL

---

## Batas tanggung jawab (penting)

| Data / aksi | Siapa yang decide |
|---|---|
| Visual, animasi, camera | Client |
| “Saya menekan tombol beli” | Client (request) |
| Apakah beli berhasil, harga, stok, saldo baru | Backend (+ game server orchestration) |
| Possession kendaraan | Game server |
| Saldo bank / inventory slots | Backend + DB |
| Cuaca / time of day | Game server (replicated) |

**Rule:** Client tidak pernah “menetapkan” angka uang atau isi inventory.

---

## Networking model (MVP)

- Dedicated server (bukan listen-server pemain sebagai host)
- Server authoritative movement dengan client-side prediction + reconciliation (standar UE)
- Chat via game server atau backend WS — pilih satu path di implementasi awal, dokumentasikan di sini saat dipilih
- Tidak ada trust pada RPCs yang mengubah ekonomi tanpa round-trip backend

### Alur contoh: beli item di warung

1. Client kirim request `BuyItem(itemId, qty)` ke game server
2. Game server cek jarak/interaksi valid
3. Game server → Backend: `POST /economy/purchase` (characterId, itemId, qty, shopId)
4. Backend: transaksi DB (lock/saldo/inventory), audit log
5. Backend return hasil → game server → client update UI

---

## Auth & session

1. Client login ke Backend API → terima JWT
2. Client connect ke Game Server dengan token
3. Game Server validasi token ke Backend (atau verify JWT + session Redis)
4. Character load dari DB; spawn di world
5. Disconnect → flush state penting (posisi, inventory dirty flags, vehicle)

---

## Skema data (konseptual MVP)

Tabel inti (nama final mengikuti migrasi):

- `accounts`
- `characters`
- `inventories` / `inventory_items`
- `wallets` (cash)
- `bank_accounts`
- `transactions` (audit)
- `vehicles`
- `houses`
- `sessions` (opsional di Redis saja untuk ephemeral)

Detail kolom → dibuat saat implementasi migrasi (`backend/migrations/`).

---

## Keamanan (baseline)

- Server-side validation semua aksi penting
- Rate limiting (API + game actions sensitif)
- TLS untuk API; channel game sesuai setup UE
- Encryption at rest sesuai provider DB/storage
- Anti-cheat dasar: reject impossible values, speed sanity, distance checks
- Logging + audit trail untuk setiap mutasi ekonomi

---

## Scalability path (jangan overbuild di hari 1)

**Sekarang (MVP):** 1 game server instance + 1 API + 1 DB + 1 Redis.

**Nanti:**
- Multiple dedicated servers (shards / instances)
- API horizontal scale behind load balancer
- Read replicas DB jika perlu
- Redis cluster
- Matchmaking service memetakan player → instance

Arsitektur MVP harus **tidak menghalangi** langkah di atas (characterId, serverId, instance routing sejak awal di desain session).

---

## Repo layout (disarankan)

```
aldopr/
├── README.md
├── docs/                 # visi, MVP, arsitektur, roadmap
├── backend/              # Go API, migrations, workers
├── game/                 # Unreal project (atau submodule)
├── infra/                # docker-compose, scripts deploy
└── tools/                # admin scripts, codegen, etc.
```

Folder `game/` dapat diisi setelah Unreal project dibuat; docs & backend boleh lebih dulu.

---

## Referensi

- [MVP](./MVP.md)
- [Tech Decisions](./TECH_DECISIONS.md)
- [Roadmap](./ROADMAP.md)
