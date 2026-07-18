# MVP 0.1 — Scope Lock

Dokumen ini **mengunci** apa yang masuk Versi 0.1. Fitur di luar daftar ini dianggap **out of scope** sampai milestone berikutnya.

Target: fondasi yang bisa dimainkan bersama (20–50 pemain), bukan konten sebanyak mungkin.

---

## Definisi sukses MVP

Pemain dapat:

1. Login / register
2. Membuat karakter
3. Masuk ke open world kecil bersama pemain lain
4. Chat
5. Naik motor & mobil
6. Punya inventory + uang
7. Pakai bank & SPBU
8. Interaksi warung & rumah sederhana
9. Progress tersimpan (database)
10. Melihat NPC dasar + siang/malam + cuaca

Jika semua di atas stabil, MVP **lulus**. Sisanya ditunda.

---

## In scope (wajib)

### Akun & karakter
- [x] Register / login (JWT) — API backend scaffold
- [x] Character creation dasar (gender, kulit, rambut, wajah, pakaian starter) — API fields
- [x] 1 karakter per akun (multi-character bisa belakangan)
- [ ] Integrasi flow ini dari UE client / dedicated server

### Dunia & multiplayer
- [ ] Open world **kecil** (1–2 zona saja, mis. Downtown + Perumahan/Kampung ringkas)
- [ ] Dedicated server, server-authoritative movement & interaksi penting
- [ ] 20–50 concurrent players (target desain; optimasi bertahap)
- [ ] Chat text (global / proximity — pilih satu untuk 0.1, jangan keduanya kalau belum siap)

### Kendaraan
- [x] Katalog motor + mobil (API, desain fiksi)
- [x] Ownership + kunci di inventory
- [x] Bensin + isi di SPBU (API refuel)
- [ ] Enter/exit, drive dasar (UE)
- [ ] Lampu / klakson (nice-to-have jika waktu cukup; tidak memblokir MVP)

### Jobs (loop uang sederhana)
- [x] 3 job shift server-side (ojol/kurir/kasir) — payout dikunci server
- [ ] Integrasi interact job di UE

### Ekonomi & inventory
- [x] Dompet (cash on hand) — field + starting cash
- [x] Bank (saldo + deposit/withdraw sederhana) — API
- [x] Inventory server-side (makanan/minuman/tools dasar)
- [x] Transaksi warung (beli item)
- [x] Semua mutasi uang/item **hanya di server** (pola bank sudah)

### Properti
- [x] Rumah sederhana: beli (ownership) — API
- [x] Interior minimal / spawn point rumah — pos_x/y/z di listing
- [x] Kunci rumah di inventory (`house_key_*`)
- [ ] Integrasi enter/exit rumah di UE

### Dunia hidup (minimal)
- [ ] Siklus siang/malam
- [ ] Cuaca dinamis dasar (sunny / cloudy / rain)
- [ ] NPC dasar: walk / idle / simple schedule (Behavior Tree), **tanpa** LLM

### Persistensi & infrastruktur
- [x] PostgreSQL: account, character, inventory, money, house (+ position save)
- [x] Save progress posisi karakter (API `PATCH .../position`)
- [ ] Save progress saat disconnect / interval aman (UE hook)
- [ ] Redis: session / online presence (minimal)
- [x] Logging transaksi dasar (tabel `transactions`)
- [x] Transfer cash antar pemain (API)
- [x] Job cooldown anti-spam (60 detik / job)

---

## Out of scope (sengaja ditunda)

Jangan kerjakan ini di 0.1 kecuali sisa kapasitas setelah in-scope selesai:

| Area | Ditunda ke |
|---|---|
| Bisnis player (warung milik pemain, offline income) | 0.2+ |
| Marketplace / trading player-to-player | 0.2+ |
| ATM UI lengkap, transfer antar pemain | 0.2 |
| Banyak pekerjaan (ojol, polisi, dll.) — cukup 1–2 loop uang sederhana | 0.2 |
| Furniture move & housing decoration dalam | 0.2–0.3 |
| Bus, truk, kapal, sepeda | 0.3 |
| Guild, friend list UI lengkap | 0.3 |
| Quest/achievement system penuh | 0.2–0.3 |
| Pajak, harga dinamis kompleks | 0.3+ |
| Crafting | later |
| Voice chat | later |
| Anti-cheat lanjutan | bertahap setelah 0.1 |
| Map zona lengkap (pantai, bandara, sawah, dll.) | bertahap |
| AAA art final / seluruh asset original | bertahap (pakai Megascans/Fab dulu) |
| LLM untuk NPC | bukan prioritas awal |

---

## Batasan teknis MVP

- **1 region / 1 shard** dulu (belum multi-region matchmaking kompleks)
- Client **tidak** boleh menentukan saldo, inventory, ownership
- Asset: gratis / berlisensi (Quixel, Fab, Unreal free). Placeholder OK.
- Target visual: “bisa dibanggakan”, bukan “AAA final”
- FPS: usahakan 60 di mid-range; profile early, jangan tunda optimasi total ke akhir

---

## Kriteria “Done” per fitur

Sebuah fitur MVP dianggap selesai hanya jika:

1. Jalan di dedicated server (bukan hanya PIE single player)
2. Validasi server-side ada
3. Data penting tersimpan di DB
4. Ada cara reproduksi / checklist uji singkat di `docs/` atau test plan internal

---

## Anti-pattern (hindari)

- Memperbesar map sebelum movement + sync + save stabil
- Menambah 10 pekerjaan sebelum ekonomi cash/bank aman
- Customizer karakter ultra-detail sebelum login → world → save flow lengkap
- Menyalin merek/kendaraan/logo nyata

---

## Referensi

- [Vision](./VISION.md)
- [Architecture](./ARCHITECTURE.md)
- [Roadmap](./ROADMAP.md)
