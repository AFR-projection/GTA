# Indonesia Life Online

Open World Multiplayer Life Simulation beridentitas Indonesia — dibangun dengan **Unreal Engine 5**, dedicated server, dan backend server-authoritative.

> Working title. Bukan GTA clone: tujuannya dunia virtual Indonesia yang **terasa hidup**.

---

## Status

| Item | Status |
|---|---|
| Fase | **0 — Project foundation** |
| MVP target | **0.1** (lihat docs) |
| Engine | Unreal Engine 5 |
| Backend | Go + PostgreSQL + Redis |
| Repo | Docs foundation siap; game & backend scaffold menyusul |

---

## Visi singkat

Pemain login dan langsung merasa: **"Ini Indonesia banget."**

Hidup, kerja, bisnis, rumah, kendaraan, sosialisasi — cerita dibuat oleh pemain, bukan campaign linear. Detail suasana (warung, gang, sawah, tol, pasar, dll.) lebih penting daripada map raksasa kosong.

Dokumen lengkap: [`docs/VISION.md`](docs/VISION.md)

---

## Mulai dari sini (untuk kontributor)

Baca berurutan:

1. [`docs/VISION.md`](docs/VISION.md) — kenapa proyek ini ada  
2. [`docs/MVP.md`](docs/MVP.md) — **apa yang boleh dikerjakan di 0.1** (scope lock)  
3. [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — bagaimana sistem saling bicara  
4. [`docs/TECH_DECISIONS.md`](docs/TECH_DECISIONS.md) — stack yang dikunci & alasan  
5. [`docs/ROADMAP.md`](docs/ROADMAP.md) — urutan versi  

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

## Struktur repo (target)

```
aldopr/
├── README.md          ← kamu di sini
├── docs/              ← visi, MVP, arsitektur, roadmap, keputusan teknis
├── backend/           ← Go API, migrations (menyusul)
├── game/              ← Unreal project (menyusul)
├── infra/             ← docker-compose, deploy scripts (menyusul)
└── tools/             ← utilitas tim (menyusul)
```

---

## MVP 0.1 (ringkas)

Login, character creation, open world kecil, multiplayer, chat, motor + mobil, inventory, uang, bank, SPBU, warung, rumah sederhana, save ke database, NPC dasar, siang/malam, cuaca. Target **20–50 pemain** / server.

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

Setup lokal backend/game akan ditambahkan setelah scaffold Fase 0 selesai.

---

## Lisensi

TBD — tentukan sebelum distribusi publik / playtest eksternal.
