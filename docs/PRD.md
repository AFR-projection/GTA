# Product Requirements Document (PRD)

**Product:** Indonesia Life Online (working title)  
**Document owner:** Product / Tech Lead  
**Status:** Living document — update saat scope berubah  
**Last updated:** 2026-07-19  
**Related:** [VISION](./VISION.md) · [MVP 0.1](./MVP.md) · [ARCHITECTURE](./ARCHITECTURE.md) · [ROADMAP](./ROADMAP.md) · [TECH_DECISIONS](./TECH_DECISIONS.md)

---

## 1. Ringkasan produk

Indonesia Life Online adalah game **Open World Multiplayer Life Simulation** di Windows PC (Unreal Engine 5) yang menempatkan pemain di dunia virtual beridentitas Indonesia.

Pemain menjalani kehidupan sehari-hari bersama pemain lain: bergerak di kota fiksi, bekerja, berbelanja, punya uang & inventory, membeli rumah/kendaraan, dan bersosialisasi — **tanpa** campaign utama yang memaksa cerita.

Ini **bukan** GTA clone. Keberhasilan diukur dari rasa “hidup di Indonesia”, bukan dari map terbesar atau combat.

---

## 2. Problem statement

| Masalah | Dampak |
|---|---|
| Sedikit game multiplayer life-sim yang terasa lokal Indonesia | Pemain Indonesia jarang dapat world yang familiar secara budaya & visual |
| Open world sering fokus luas tapi kosong | Dunia terasa mati; retensi rendah |
| Banyak proyek gagal karena scope meledak + client-authoritative | Cheat, ekonomi rusak, sulit scale |

**Solusi kami:** dunia fiksi bergaya Indonesia yang terasa hidup, fondasi **server-authoritative**, rilis bertahap (MVP → live service).

---

## 3. Goals & non-goals

### Goals (produk)
1. Pemain login dan langsung merasa: **“Ini Indonesia banget.”**
2. Pemain bisa menciptakan cerita sendiri lewat aktivitas harian + interaksi sosial.
3. Sistem uang, inventory, rumah, kendaraan **aman** (server-side).
4. Fondasi teknis siap berkembang bertahun-tahun lewat update.

### Non-goals (bukan fokus awal)
- Replika Indonesia 1:1 / map raksasa kosong
- Story campaign linear ala single-player GTA
- Combat/PvP sebagai pilar utama
- LLM untuk NPC di tahap awal
- Mobile / konsol (fase 0.1)

---

## 4. Target pengguna

### Persona utama
- **Pemain life-sim / sandbox online** (PC), usia ~16–30
- Suka sosialisasi, roleplay ringan, ekonomi, kendaraan, rumah
- Familiar dengan nuansa kota Indonesia (warung, gang, motor, perumahan, dll.)

### Persona sekunder
- Content creator / streamer yang butuh world unik lokal
- Pemain Roblox/GTA RP yang ingin alternatif berbahasa & bersuasana Indonesia

### Bukan target awal
- Hardcore competitive esports
- Simulasi militer / racing murni

---

## 5. Value proposition

> Dunia virtual Indonesia yang hidup — kerja, belanja, rumah, kendaraan, nongkrong bareng teman — dibangun untuk tumbuh lama, bukan demo sekali main.

Pembeda:
- Identitas visual & audio Indonesia (bukan reskin kota barat)
- Life simulation + multiplayer, bukan shooter
- Ekonomi & ownership server-side sejak hari 1

---

## 6. Produk experience (user journeys)

### Journey A — First session (MVP)
1. Register / login  
2. Buat karakter  
3. Spawn di open world kecil bersama pemain lain  
4. Jalan / chat  
5. Belanja warung → inventory bertambah, cash berkurang  
6. Deposit/withdraw bank  
7. Beli rumah sederhana → dapat kunci  
8. Disconnect → progress tersimpan  

### Journey B — Daily loop (pasca-MVP)
Login → kerja/job loop → hasilkan uang → upgrade rumah/kendaraan → sosialisasi → logout aman.

### Journey C — Social
Bertemu pemain lain, chat, (nanti) grup/friend, aktivitas bersama di dunia.

---

## 7. Functional requirements

Prioritas: **P0** wajib MVP 0.1 · **P1** penting segera setelah · **P2** later.

### 7.1 Account & character
| ID | Requirement | Priority | Status |
|---|---|---|---|
| FR-A1 | Register & login dengan kredensial aman | P0 | Done (API) |
| FR-A2 | Session via JWT untuk API | P0 | Done |
| FR-A3 | Character creation modular dasar | P0 | Done (API fields) |
| FR-A4 | 1 karakter / akun di MVP | P0 | Done |
| FR-A5 | UE client memakai flow auth/character | P0 | Todo |

### 7.2 World & multiplayer
| ID | Requirement | Priority | Status |
|---|---|---|---|
| FR-W1 | Open world kecil (1–2 zona) | P0 | Todo (UE) |
| FR-W2 | Dedicated server, server-authoritative | P0 | Todo (UE) |
| FR-W3 | Target desain 20–50 concurrent / instance | P0 | Todo |
| FR-W4 | Global text chat | P0 | Todo |
| FR-W5 | Day/night cycle | P0 | Todo |
| FR-W6 | Weather dasar (sunny/cloudy/rain) | P0 | Todo |
| FR-W7 | NPC dasar (BT/State Machine, no LLM) | P0 | Todo |

### 7.3 Economy & inventory
| ID | Requirement | Priority | Status |
|---|---|---|---|
| FR-E1 | Cash on hand server-side | P0 | Done |
| FR-E2 | Bank deposit/withdraw + audit | P0 | Done |
| FR-E3 | Inventory persistent | P0 | Done |
| FR-E4 | Belanja warung (harga dari server) | P0 | Done |
| FR-E5 | Client tidak boleh set saldo/item | P0 | Done (pola) |
| FR-E6 | Marketplace / transfer P2P | P1 | Later (0.2) |
| FR-E7 | Dynamic pricing / pajak | P2 | Later |

### 7.4 Housing
| ID | Requirement | Priority | Status |
|---|---|---|---|
| FR-H1 | Beli rumah (ownership) | P0 | Done (API) |
| FR-H2 | Spawn point / posisi rumah tersimpan | P0 | Done (API) |
| FR-H3 | Kunci rumah di inventory | P0 | Done (API) |
| FR-H4 | Enter/exit & interior di UE | P0 | Todo |
| FR-H5 | Furniture / decorate | P1 | Later |

### 7.5 Vehicles
| ID | Requirement | Priority | Status |
|---|---|---|---|
| FR-V1 | 1 motor + 1 mobil prototype | P0 | Todo |
| FR-V2 | Enter/exit + drive dasar | P0 | Todo |
| FR-V3 | Fuel + SPBU | P0 | Todo |
| FR-V4 | Ownership persistent | P0 | Todo |
| FR-V5 | Damage/repair dalam | P1 | Later |

### 7.6 Jobs & business
| ID | Requirement | Priority | Status |
|---|---|---|---|
| FR-J1 | 1–2 loop uang sederhana | P1 | 0.2 |
| FR-J2 | Banyak job (ojol, dll.) | P1–P2 | 0.2+ |
| FR-J3 | Player-owned business + offline income | P1 | 0.2+ |

### 7.7 Platform & ops
| ID | Requirement | Priority | Status |
|---|---|---|---|
| FR-O1 | Persistensi PostgreSQL (Neon) | P0 | In progress |
| FR-O2 | Redis session/presence | P1 | Optional sekarang |
| FR-O3 | Object storage (avatar/screenshot) | P2 | Later |
| FR-O4 | Rate limit + audit transaksi | P0 | Partial (audit ada) |
| FR-O5 | Anti-cheat dasar (sanity checks) | P1 | Bertahap |

Detail checklist teknis: [MVP.md](./MVP.md).

---

## 8. Non-functional requirements

| ID | Area | Requirement |
|---|---|---|
| NFR-1 | Performance | Target 60 FPS di hardware menengah (PC) |
| NFR-2 | Authority | Semua mutasi ekonomi/inventory/ownership di server/backend |
| NFR-3 | Scalability | Arsitektur MVP tidak menghalangi multi-instance nanti |
| NFR-4 | Security | Secret di env; JWT; jangan trust client |
| NFR-5 | Legal / IP | Desain orisinal; larang salin merek/logo nyata |
| NFR-6 | Maintainability | Modular, documented, conventional commits |
| NFR-7 | Localization | Bahasa UI awal: Indonesia (EN optional later) |

---

## 9. Success metrics

### MVP 0.1 (internal)
- 2+ client connect dedicated server tanpa crash rutin
- Register → character → warung → rumah → relog data utuh
- 0 critical economy exploit yang diketahui di playtest internal
- Playtest 20+ concurrent (stress) tanpa hard lock economy

### Pasca-MVP (arah)
- Session length & D1/D7 retention (ukur setelah soft launch)
- % pemain yang beli rumah / kendaraan dalam 3 sesi pertama
- Ticket cheat ekonomi / minggu (harus turun seiring hardening)

---

## 10. Constraints & assumptions

**Constraints**
- Platform: Windows PC dulu
- Engine: Unreal Engine 5
- Backend: Go + PostgreSQL (Neon) + Redis (nanti)
- Tim kecil → scope harus kejam

**Assumptions**
- Asset awal dari Megascans/Fab/gratis berlisensi
- Docker lokal tidak wajib; Neon cukup untuk dev
- Nama produk final bisa berubah

---

## 11. Risks

| Risk | Mitigation |
|---|---|
| Scope creep open world | MVP.md lock; PRD non-goals |
| UE networking sulit | Slice kecil dulu: movement + chat sebelum map besar |
| Economy exploit | Server authority + audit `transactions` |
| Art bandwidth | Placeholder dulu; swappable assets |
| Burnout “bikin GTA” | Ulangi prinsip: hidup > besar |

---

## 12. Release plan (ringkas)

| Version | Fokus produk |
|---|---|
| **0.1** | Fondasi playable: auth, world kecil, chat, kendaraan dasar, ekonomi, warung, rumah, NPC/cuaca minimal |
| **0.2** | Jobs loop, marketplace/transfer, business v1, furniture |
| **0.3** | World expansion, lebih banyak kendaraan, NPC lebih hidup |
| **0.4+** | Live ops, moderation, multi-instance, polish art |

Detail: [ROADMAP.md](./ROADMAP.md).

---

## 13. Hierarchy dokumen

```
PRD.md              ← apa yang dibangun & kenapa (produk)
  ├── VISION.md     ← arah kreatif jangka panjang
  ├── MVP.md        ← scope lock versi 0.1 (checklist)
  ├── ROADMAP.md    ← urutan rilis
  ├── ARCHITECTURE.md / TECH_DECISIONS.md  ← bagaimana teknis
  └── ONBOARDING.md ← tugas kontributor
```

**Aturan:** ubah fitur produk → update PRD + MVP. Ubah cara implementasi saja → Architecture/Tech Decisions.

---

## 14. Open product questions

| Question | Current decision | Revisit |
|---|---|---|
| Chat model 0.1 | Global | 0.2 proximity |
| Housing 0.1 | Buy ownership, 1 house/char | Multi-house later |
| Monetization | TBD (jangan desain pay-to-win) | sebelum public launch |
| Final product name | Working title ILO | sebelum marketing |

---

## 15. Approval

| Role | Name | Date | Sign-off |
|---|---|---|---|
| Product / Owner | | | |
| Tech Lead | | | |

*Isi saat tim formalized. Untuk sekarang dokumen ini adalah sumber kebenaran produk di repo.*
