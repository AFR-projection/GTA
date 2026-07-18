# Onboarding — tugas konkret

Anggap lead sudah set arah. Ini yang **kamu** kerjakan biar proyek maju, bukan stuck di fantasi.

---

## Hari ini (blocker mesin)

Tanpa ini, backend tidak bisa dijalankan:

1. Install **Go** — https://go.dev/dl/  
2. Install **Docker Desktop** — https://www.docker.com/products/docker-desktop/  
3. Restart terminal / PC jika PATH belum ke-update  
4. Dari root repo:

```powershell
Copy-Item .env.example .env
cd infra
docker compose up -d
cd ..\backend
go mod tidy
go run .\cmd\api
```

5. Buka `http://localhost:8080/healthz` → harus `{"status":"ok"}`  
6. Jalankan smoke test di [`backend/README.md`](../backend/README.md)

Laporkan ke lead kalau stuck di install (screenshot error).

---

## Minggu ini (urutan kerja — jangan loncat)

| Prioritas | Task | Done when |
|---|---|---|
| P0 | API auth + character jalan lokal | register → login → create character OK |
| P0 | Migrasi DB idempotent | restart API tidak rusak schema |
| P1 | Economy stub endpoints | deposit/withdraw bank + audit `transactions` |
| P1 | Inventory add/remove server-side | item warung bisa dibeli (tanpa UE dulu) |
| P2 | Buat Unreal project kosong di `game/` | 2 PIE clients connect dedicated server (tanpa backend dulu juga OK) |

**Dilarang minggu ini:** bikin map besar, customizer ultra, 10 pekerjaan, bisnis offline income.

---

## Keputusan yang sudah dikunci (jangan debat lagi)

- Chat 0.1 = **global**
- Rumah 0.1 = **beli**
- 1 karakter / akun
- Backend = **Go**

Lihat [`TECH_DECISIONS.md`](./TECH_DECISIONS.md).

---

## Cara minta review ke lead

Kirim:

1. Apa yang berubah (1–3 bullet)
2. Cara repro / test
3. Link commit / PR
4. Apa yang masih rusak / belum

Jangan kirim “udah progress” tanpa bukti jalan.
