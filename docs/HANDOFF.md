# HANDOFF — baca ini dulu (untuk AI agent / developer baru)

**Product:** Indonesia Life Online (working title)  
**Repo:** https://github.com/AFR-projection/GTA  
**Local path:** `C:\Users\User\Documents\aldopr`  
**Last updated:** 2026-07-19  

Kalau kamu AI agent (Cursor, Claude Code, dll.): **jangan mulai dari nol.** Baca urutan di bawah, lalu lanjut task di bagian “NEXT”.

---

## 1) Baca dokumen ini dulu (urut)

1. `docs/HANDOFF.md` ← kamu di sini  
2. `docs/PROGRESS.md` ← checklist fase (sumber kebenaran “sudah sampai mana”)  
3. `docs/PRD.md` ← requirements produk  
4. `docs/MVP.md` ← scope lock 0.1  
5. `docs/ARCHITECTURE.md` + `docs/TECH_DECISIONS.md`  
6. `docs/UE_SETUP.md` + `docs/UE_INTEGRATION.md` (kalau kerjakan Unreal)  
7. `backend/README.md` ← endpoint API  

---

## 2) Stack & keputusan terkunci

| Item | Keputusan |
|---|---|
| Engine | Unreal Engine 5 (**5.8 OK** — user sedang install) |
| Backend | **Go** API, server-authoritative |
| DB | **Neon Postgres** (bukan Docker lokal wajib) |
| Cache | Redis — belum dipakai API |
| Chat 0.1 | Global |
| Rumah 0.1 | Beli ownership, **1 rumah / karakter** |
| Character | **1 karakter / akun** |
| Job cooldown | 60 detik / job_key |

**Jangan** usulkan ganti stack tanpa update `TECH_DECISIONS.md`.

---

## 3) Cara jalanin lokal (Windows PowerShell)

```powershell
# API (pastikan port 8080 kosong)
cd C:\Users\User\Documents\aldopr\backend
go run .\cmd\api
```

- Env: root `.env` (dari `.env.example`) — `DATABASE_URL` = Neon URI  
- Health: `http://localhost:8080/healthz`  
- PowerShell: pakai `Invoke-RestMethod`, **bukan** `curl -H` (alias beda)

Test accounts yang sudah ada di Neon (dev):

| Email | Password | Character |
|---|---|---|
| `dev@ilo.local` | `password123` | Budi |
| `dev2@ilo.local` | `password123` | Siti (`e0190a81-6373-4815-bcac-3a9dee936471`) |

**Jangan commit** `.env` / secrets.

---

## 4) Apa yang SUDAH hidup (backend)

- Auth JWT: register / login / me  
- Characters CRUD (max 1) + `PATCH .../position` + `GET .../summary`  
- Bank deposit/withdraw (audit `transactions`)  
- Warung catalog + purchase → inventory  
- Housing listings + buy (1 house) + house key item  
- Vehicles catalog + buy + refuel SPBU + consume-fuel + vehicle position  
- Jobs (ojol/kurir/kasir) + cooldown table  
- P2P cash transfer  

Migrations goose: `00001` … `00004` (current version **4**).

---

## 5) Apa yang BELUM (blocker / next)

### Blocker visual
- Unreal project **belum** ada di `game/` (tunggu UE 5.8 selesai install)  
- Belum: dedicated server, movement sync, chat in-world, UE↔API login UI  

### Backend nice-to-have (boleh dikerjakan sambil nunggu UE)
- Rate limiting HTTP  
- Consume/use inventory items (makan)  
- Vehicle ownership list polish / sell vehicle  
- Redis session (opsional)  

### Jangan kerjakan dulu
- Map raksasa, customizer ultra, 10 job kompleks, bisnis offline income, LLM NPC  

---

## 6) Prompt yang bisa ditempel ke AI agent lain

Copy-paste:

```text
Lanjutkan project Indonesia Life Online di repo ini.

WAJIB baca dulu:
- docs/HANDOFF.md
- docs/PROGRESS.md
- docs/PRD.md
- docs/MVP.md

Ikuti TECH_DECISIONS yang sudah dikunci.
Update docs/PROGRESS.md setiap selesai satu fase/task (centang + tanggal).
Jangan commit .env.
Conventional commits (feat/fix/docs/chore).
Remote: origin = https://github.com/AFR-projection/GTA.git (branch main).

NEXT prioritas:
1) Jika UE 5.8 sudah terinstall → buat project ILO di folder game/ (Third Person), ikuti docs/UE_SETUP.md + UE_INTEGRATION.md
2) Jika UE belum siap → lanjut backend dari PROGRESS.md bagian "Next up"
```

---

## 7) Aturan kerja untuk agent

1. Satu fokus per sesi (jangan loncat map + economy sekaligus).  
2. Setelah fitur selesai: update **`docs/PROGRESS.md`** (checklist).  
3. Push ke `main` hanya jika user minta / sudah jadi kebiasaan repo ini.  
4. Test API dengan PowerShell `Invoke-RestMethod`.  
5. Economy selalu server-side.  

---

## 8) Kontak konteks user

- User = owner / “anak buah” belajar; agent = lead teknis.  
- Bahasa: Indonesia, santai tapi profesional.  
- User **tidak mau Docker lokal** — pakai Neon.  
