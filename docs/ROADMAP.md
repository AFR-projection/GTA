# Roadmap

Roadmap ini mengikuti prinsip: **fondasi dulu, konten kemudian**. Setiap versi harus menambah nilai bermain tanpa merusak arsitektur inti.

---

## Fase 0 — Project foundation (sekarang)

**Tujuan:** tim bisa kerja profesional di repo yang sama.

- [x] Vision document
- [x] MVP scope lock
- [x] Architecture + tech decisions
- [x] README proyek
- [x] Struktur folder repo (`backend/`, `infra/`, `game/`, `tools/`)
- [x] Git init + `.gitignore` (UE + Go + secrets)
- [x] Docker Compose: Postgres + Redis lokal (**optional** — default tim pakai Neon)
- [x] Backend scaffold: auth JWT + character CRUD + migrations
- [x] Coding standards singkat (commit convention di TECH_DECISIONS)
- [ ] Contributor: Neon `DATABASE_URL` terisi & API `healthz` hijau

**Exit:** developer baru bisa clone → isi Neon URL → `go run` → register/login/character.

---

## Versi 0.1 — Playable multiplayer slice

Lihat detail wajib di [MVP.md](./MVP.md).

Urutan build yang disarankan (jangan dibalik tanpa alasan kuat):

1. **Backend auth + character CRUD** + migrasi DB  
2. **UE dedicated server empty level** + connect dengan auth  
3. **Movement sync** multiplayer  
4. **Chat**  
5. **Economy stub** (wallet/bank) + inventory  
6. **Warung purchase flow** (end-to-end server-side)  
7. **Kendaraan** (motor + mobil) + SPBU  
8. **Rumah sederhana**  
9. **Day/night + weather**  
10. **NPC dasar**  
11. **Save/load harden** + stress 20+ players  

**Exit:** checklist MVP.md semua in-scope tercentang dan stabil di dedicated server.

---

## Versi 0.2 — Economy & jobs loop

- Transfer antar pemain / ATM UX
- 2–4 pekerjaan dengan loop jelas (contoh: ojol, kurir, kasir, petani sederhana)
- Marketplace dasar atau trading aman
- Business ownership v1 (satu tipe: warung/cafe) — income rules tertulis
- Housing: interior furniture place/move dasar
- Friend list dasar
- Quest framework sederhana

**Exit:** pemain punya alasan login harian (uang + progres pekerjaan + rumah).

---

## Versi 0.3 — World expansion & vehicles

- Zona baru: pasar, pantai, atau kampung/sawah (pilih berdasarkan art bandwidth)
- Bus / truk / sepeda (prioritas sesuai konten)
- NPC density & schedule lebih kaya
- Damage/repair kendaraan lebih dalam
- Guild / group dasar
- Optimasi networking & streaming level

**Exit:** dunia terasa lebih “Indonesia” dan lebih ramai tanpa merusak FPS target.

---

## Versi 0.4+ — Live service foundation

- More businesses, dynamic prices, pajak (jika masih diinginkan)
- Events, season updates
- Anti-cheat lanjutan, moderation tools
- Replay/screenshot pipeline (R2)
- Admin dashboard
- Possibly multiple server instances + matchmaking

Konten art original menggantikan placeholder secara bertahap (**asset swappable** tanpa rewrite gameplay).

---

## Prinsip prioritas tiap sprint

Jika konflik waktu, urutkan:

1. Stability & correctness (ekonomi/inventory/save)
2. Multiplayer feel (latency, desync)
3. Core loop fun (gerak, kendaraan, belanja, rumah)
4. World living (NPC, cuaca, audio)
5. Visual polish

Grafik AAA final **tidak** boleh menahan 0.1.

---

## Milestone review

Setiap akhir versi:

1. Main bareng (internal playtest)
2. Catat bug P0/P1
3. Update MVP/roadmap jika scope bergeser (ubah dokumen, jangan diam-diam)
4. Tag release di git (`v0.1.0`, dst.)

---

## Referensi

- [Vision](./VISION.md)
- [MVP](./MVP.md)
- [Architecture](./ARCHITECTURE.md)
