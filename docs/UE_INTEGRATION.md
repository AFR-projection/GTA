# UE ↔ Backend Integration Contract

Kontrak ini supaya Unreal dan Go API tidak saling tebak.

**Backend base URL (dev):** `http://127.0.0.1:8080`  
**Auth:** `Authorization: Bearer <jwt>`

---

## Urutan integrasi (wajib ikut)

| Step | Di UE | Backend |
|---|---|---|
| 0 | Dedicated server + 2 clients movement | — |
| 1 | UI Login/Register → HTTP | `/v1/auth/*` |
| 2 | Character create/select | `/v1/characters` |
| 3 | Setelah spawn, periodik save posisi | `PATCH /v1/characters/{id}/position` |
| 4 | Warung interact → purchase | `/v1/characters/{id}/shops/warung/purchase` |
| 5 | Housing interact → buy / list | `/v1/housing/*`, `/houses/buy` |
| 6 | Bank ATM UI | deposit/withdraw |
| 7 | Vehicle possession + fuel | (API vehicle menyusul) |

**Jangan** mulai step 4–7 sebelum step 0–2 stabil.

---

## HTTP dari UE

Gunakan `FHttpModule` (C++) atau plugin HTTP Blueprint.

### Register
`POST /v1/auth/register`

```json
{ "email": "a@b.com", "password": "password123", "display_name": "Nama" }
```

### Login
`POST /v1/auth/login`

```json
{ "email": "a@b.com", "password": "password123" }
```

Response: simpan `token` di GameInstance (jangan hardcode di Blueprint asset yang di-commit).

### Create character
`POST /v1/characters` + Bearer

```json
{
  "name": "Siti",
  "gender": "female",
  "skin_tone": 2,
  "hair_style": 1,
  "face_preset": 0,
  "outfit_id": "starter_01"
}
```

### Load character (full hydrate)
`GET /v1/characters/{id}/summary` → character + inventory + houses + vehicles (pakai ini saat login world).

### Save position (server truth untuk disconnect)
`PATCH /v1/characters/{id}/position` + Bearer

```json
{ "pos_x": 10.5, "pos_y": 0.0, "pos_z": -3.2 }
```

### Vehicle while driving (UE dedicated server / interval)
- `POST .../vehicles/{vehicleID}/consume-fuel` `{ "amount": 0.5 }`
- `PATCH .../vehicles/{vehicleID}/position` `{ "pos_x", "pos_y", "pos_z" }`

### Load character (single)
`GET /v1/characters/{id}` → pakai `pos_*`, `cash`, `bank` untuk hydrate state UI.

---

## Authority rules (wajib)

| Data | Siapa yang decide |
|---|---|
| Visual movement prediksi | Client |
| Final possession / interact jarak | UE Dedicated Server |
| Cash, inventory, house ownership | **Go API only** |
| Harga warung / listing rumah | **Go catalog only** |

Dedicated server UE memanggil Go API (atau client kirim intent → DS validasi jarak → DS panggil API).  
**Jangan** biarkan client langsung “saya set cash = 999999”.

Recommended MVP flow belanja:

```
Client interact Warung
  → Server RPC BuyItem(itemKey, qty)
  → DS cek jarak ke shop actor
  → DS HTTP POST purchase ke Go (dengan service token / player JWT)
  → DS apply hasil ke player state / UI
```

Untuk minggu 1 boleh **client HTTP langsung** ke Go (lebih cepat prototyping), lalu pindah ke DS-mediated sebelum playtest publik.

---

## Chat 0.1

- **Global text chat** (keputusan terkunci)
- Boleh murni UE replication dulu (tidak wajib lewat Go)
- Persist chat history = bukan MVP

---

## Map 0.1

- 1 level kecil: plaza/kampung strip + warung + marker rumah + SPBU placeholder
- Streaming besar = later

---

## Definition of Done — UE slice pertama

- [ ] `ILO.uproject` ada di `game/` dan bisa dibuka Editor  
- [ ] Dedicated server package atau `play as listen` dilarang untuk target; usahakan dedicated  
- [ ] 2 editor/clients: saling lihat movement  
- [ ] Chat global kelihatan  
- [ ] Login UI → token → create/load character dari API  
- [ ] Posisi tersimpan setelah relog  

---

## Referensi API

Lihat [`backend/README.md`](../backend/README.md) untuk daftar endpoint terkini.
