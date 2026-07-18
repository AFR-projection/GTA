# Progress Checklist — Indonesia Life Online

**Cara pakai:** setiap task selesai → centang `[x]` + isi tanggal di kolom Notes.  
Ini file yang dibaca AI agent lain untuk tahu **sudah sampai mana**.

Update terakhir: **2026-07-19**

---

## Fase 0 — Project foundation

| Status | Task | Notes |
|---|---|---|
| [x] | Vision / PRD / MVP / Architecture / Roadmap / Tech decisions | 2026-07-19 |
| [x] | Repo GitHub + README | https://github.com/AFR-projection/GTA |
| [x] | `.gitignore` (UE + Go + secrets) | |
| [x] | HANDOFF.md untuk multi-agent | 2026-07-19 |
| [x] | Neon sebagai DB default (tanpa Docker wajib) | |
| [x] | Go terpasang di mesin user | go1.26.5 |
| [ ] | Unreal Engine terpasang + verified | User download **UE 5.8.0** (in progress) |
| [ ] | Visual Studio 2022 + C++ game workload (opsional awal) | |

---

## Fase 0.1a — Backend economy slice (API)

| Status | Task | Notes |
|---|---|---|
| [x] | Docker Compose optional di `infra/` | tidak dipakai user |
| [x] | Migrasi `00001_init` (accounts, characters, inventory, transactions) | |
| [x] | Auth register/login JWT | |
| [x] | Character create/list/get (max 1) | |
| [x] | Bank deposit/withdraw + audit | |
| [x] | Warung catalog + purchase | |
| [x] | Inventory list | |
| [x] | Housing listings + buy (1 house) + house key | migrasi `00002` |
| [x] | Vehicles buy/list + SPBU refuel | migrasi `00003` |
| [x] | Consume fuel + vehicle position | |
| [x] | Character position PATCH | |
| [x] | Character summary hydrate endpoint | |
| [x] | Jobs + cooldown 60s | migrasi `00004` |
| [x] | P2P cash transfer | |
| [x] | Neon simple-protocol fix (pgx) | |
| [ ] | HTTP rate limiting | Next up (backend) |
| [ ] | Use/consume inventory item (makan) | Next up (backend) |
| [ ] | Sell vehicle / house (opsional) | later |

---

## Fase 0.1b — Unreal multiplayer slice

| Status | Task | Notes |
|---|---|---|
| [ ] | UE 5.8 install selesai | blocker |
| [ ] | Create `game/ILO.uproject` (Third Person) di folder `game/` | lihat UE_SETUP.md |
| [ ] | 2 clients movement sync (dedicated server) | |
| [ ] | Global text chat | |
| [ ] | Login UI → JWT → load `/summary` | |
| [ ] | Save position on disconnect / interval | |
| [ ] | Warung interact → API purchase | |
| [ ] | Housing marker enter/buy flow | |
| [ ] | Spawn owned vehicle + fuel loop | |
| [ ] | Day/night + weather dasar | |
| [ ] | NPC BT dasar | |

---

## Fase 0.2+ (jangan mulai sebelum 0.1b stabil)

| Status | Task | Notes |
|---|---|---|
| [ ] | Marketplace | |
| [ ] | Player business + offline income | |
| [ ] | More jobs / loops | |
| [ ] | Furniture housing | |
| [ ] | Friend list / guild | |

---

## Next up (prioritas sekarang)

1. **Tunggu / selesaikan install UE 5.8** → centang Fase 0 UE  
2. **Buat project ILO di `game/`**  
3. Sambil nunggu UE: boleh kerjakan **rate limit** atau **use inventory item**  

---

## Cara konfirmasi ke manusia

Setelah sesi coding, agent harus bilang singkat:

- Apa yang dicentang di `PROGRESS.md`  
- Cara test  
- Apa next  

Contoh prompt user ke agent baru:

> Baca `docs/HANDOFF.md` + `docs/PROGRESS.md`, lanjut dari “Next up”. Update checklist kalau selesai.
