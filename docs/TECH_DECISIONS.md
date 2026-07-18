# Tech Decisions

Keputusan teknis yang mengikat untuk fase awal. Ubah hanya lewat diskusi tim + update dokumen ini (catat tanggal & alasan).

---

## Stack yang dikunci (v0)

| Area | Pilihan | Status |
|---|---|---|
| Game engine | Unreal Engine 5 | Locked |
| Client platform | Windows PC | Locked |
| Game networking | Dedicated Server, server-authoritative | Locked |
| Backend language | **Go** | Locked untuk MVP (revisitable setelah 0.1) |
| API style | REST (+ WebSocket bila perlu presence/notif) | Locked |
| Auth | JWT | Locked |
| Primary DB | PostgreSQL (Neon untuk cloud) | Locked |
| Cache | Redis | Locked |
| Object storage | Cloudflare R2 (atau S3-compatible) | Locked |
| NPC AI (awal) | Behavior Tree + State Machine | Locked |
| LLM untuk NPC | Tidak di tahap awal | Locked |

---

## Mengapa Go untuk backend MVP?

**Dipilih Go** karena:

1. Cocok untuk API concurrent + WebSocket ringan
2. Deploy sederhana (single binary), mudah di-container
3. Ekosistem matang untuk Postgres, Redis, JWT
4. Hiring/learning curve relatif ramah untuk backend realtime services
5. Memisahkan jelas: **UE = simulasi**, **Go = persistence & economy truth**

### Alternatif yang ditolak (untuk sekarang)

| Opsi | Alasan ditunda |
|---|---|
| .NET | Sangat valid; bisa dipertimbangkan jika tim sudah kuat C# / butuh ecosystem tertentu. Tidak dipilih agar keputusan tunggal di MVP. |
| Node.js | OK untuk prototype, kurang ideal sebagai long-term economy authority tanpa disiplin ketat |
| Logic ekonomi murni di UE server | Sulit di-scale, sulit di-audit, coupling tinggi ke game build |

Jika nanti pindah ke .NET: isolasi lewat API contract yang sama supaya game server tidak peduli implementasi.

---

## Unreal: kebijakan asset

- Prioritas: Quixel Megascans, Fab Marketplace, Unreal free/content samples
- Semua aset harus punya lisensi yang jelas untuk game komersial (cek sebelum merge)
- **Dilarang** menyalin merek, logo, desain kendaraan/produk nyata
- Placeholder diizinkan; ganti art tanpa ubah gameplay contract

---

## Networking & authority

- Dedicated server wajib untuk multiplayer nyata
- Listen-server hanya untuk eksperimen lokal, bukan arsitektur produksi
- Economy/inventory/house/vehicle ownership: Backend + DB
- Movement/interaction proximity: Game Server
- Client prediction OK; reconciliation wajib

---

## Database & migrasi

- Schema dikelola migrasi versioned di `backend/migrations/`
- Jangan ubah production schema manual tanpa migrasi
- Setiap mutasi uang → row di `transactions` (audit)

---

## Secrets & environment

- Secret tidak pernah di-commit (`.env` lokal, secret manager di cloud)
- Pisahkan config: `dev` / `staging` / `prod`
- JWT secret, DB URL, Redis URL, R2 keys = environment only

---

## Testing policy (MVP)

Minimal yang harus ada sebelum bilang “fitur selesai”:

- Backend: unit/integration untuk purchase, deposit/withdraw, inventory add/remove
- Manual playtest checklist multiplayer (2+ clients)
- Log transaksi bisa ditelusuri untuk sebuah `characterId`

Automated UE tests boleh belakangan; jangan jadi blocker 0.1, tapi backend economics harus punya tes.

---

## Naming & commit

Ikuti Conventional Commits:

- `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`

Branch suggestion:

- `main` — stabil
- `develop` — integrasi (opsional di awal; boleh `main` + feature branches dulu)
- `feat/…`, `fix/…`

---

## Keputusan terbuka (belum dikunci)

Catat di sini saat diputuskan:

| Topik | Opsi | Decide by |
|---|---|---|
| Chat: global vs proximity untuk 0.1 | global / proximity / keduanya | sebelum implementasi chat |
| Rumah 0.1: sewa vs beli | rent / buy | sebelum housing slice |
| Hosting game server | VPS / cloud VM / custom | sebelum external playtest |
| Voice chat provider | none / third-party | post-0.1 |
| Anti-cheat stack | custom rules vs middleware | post-0.1 |

---

## Change log keputusan

| Tanggal | Perubahan |
|---|---|
| 2026-07-18 | Inisialisasi dokumen; Go dipilih untuk backend MVP; stack inti dikunci |

---

## Referensi

- [Architecture](./ARCHITECTURE.md)
- [MVP](./MVP.md)
