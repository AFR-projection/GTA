# Onboarding — tugas konkret

Anggap lead sudah set arah. Ini yang **kamu** kerjakan biar proyek maju.

---

## Setup mesin (TANPA Docker lokal)

Docker lokal **tidak wajib**. Database jalan di cloud (**Neon**).

### Checklist

1. **Go** terpasang — https://go.dev/dl/ (sudah OK kalau `go version` jalan)
2. Buat project gratis di **Neon** — https://console.neon.tech  
   - Region dekat SE Asia kalau ada  
   - Copy **connection string** (URI) yang ada `sslmode=require`
3. Di root repo, pastikan ada `.env`:

```powershell
cd c:\Users\User\Documents\aldopr
Copy-Item .env.example .env
```

4. Edit `.env` — ganti `DATABASE_URL` dengan URI dari Neon.  
   **Jangan** commit / share `.env` ke chat publik.
5. Jalankan API:

```powershell
cd backend
go mod tidy
go run .\cmd\api
```

6. Cek: http://localhost:8080/healthz → `{"status":"ok"}`  
7. Smoke test register/login di [`backend/README.md`](../backend/README.md)

### Catatan

- `infra/docker-compose.yml` tetap ada buat orang yang mau DB lokal — **kamu boleh ignore**.
- Redis belum dipakai API sekarang → kosongin `REDIS_URL` boleh.

Laporkan ke lead: `healthz` hijau + 1x register berhasil.

---

## Blokir berikutnya: Unreal Engine 5

Backend slice sudah jalan. Visual multiplayer **butuh UE5**.

1. Ikuti [`UE_SETUP.md`](./UE_SETUP.md) — install Epic Launcher + UE 5.4/5.5  
2. Kabari lead: **“UE ready 5.X”**  
3. Buat project `ILO` di folder `game/`  
4. Ikuti kontrak [`UE_INTEGRATION.md`](./UE_INTEGRATION.md)

Tanpa langkah ini, kita cuma bisa lanjut API (kendaraan ownership, dll.) — bukan world yang kelihatan.

---

## Minggu ini (urutan kerja — jangan loncat)

| Prioritas | Task | Done when |
|---|---|---|
| P0 | API auth + character jalan (API lokal + Neon) | register → login → create character OK |
| P0 | Migrasi DB idempotent | restart API tidak rusak schema |
| P1 | Inventory add/remove + beli warung server-side | item bisa dibeli via API |
| P2 | Unreal project kosong di `game/` | 2 client connect dedicated server |

**Dilarang minggu ini:** map besar, customizer mega, 10 pekerjaan, bisnis offline income.

---

## Keputusan yang sudah dikunci

- Chat 0.1 = **global**
- Rumah 0.1 = **beli**
- 1 karakter / akun
- Backend = **Go**
- Dev database = **Neon cloud** (Docker lokal optional)

Lihat [`TECH_DECISIONS.md`](./TECH_DECISIONS.md).

---

## Cara minta review ke lead

1. Apa yang berubah (1–3 bullet)  
2. Cara repro / test  
3. Link commit / PR  
4. Apa yang masih rusak / belum  

Jangan kirim “udah progress” tanpa bukti jalan.
