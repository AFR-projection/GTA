# Indonesia Life Online — Vision

> Working title. Nama final dapat berubah.

## Ringkasan

**Indonesia Life Online** adalah game Open World Multiplayer Life Simulation yang dibangun dengan Unreal Engine 5. Inspirasi diambil dari GTA, The Sims, Roblox, dan life sim modern — tetapi identitas utamanya adalah **dunia virtual Indonesia yang hidup**.

Ini **bukan** GTA clone. Tujuan utamanya: pemain merasakan benar-benar hidup di Indonesia — bekerja, membangun bisnis, membeli rumah & kendaraan, bersosialisasi, dan menciptakan cerita sendiri bersama pemain lain.

Proyek ini dirancang jangka panjang, berkembang melalui update berkala selama bertahun-tahun.

---

## Visi pengalaman

Saat login, pemain harus langsung berpikir:

> "Ini Indonesia banget."

Target utama **bukan** map terbesar, melainkan dunia yang **terasa hidup**.

Detail suasana yang harus didukung (contoh):

- Gang sempit, warung kopi, warung makan, pos ronda
- Kabel listrik semrawut, spanduk pinggir jalan
- Masjid, sawah, pantai, gunung
- Jalan berlubang, jalan tol, flyover
- Perumahan subsidi, ruko, mall
- Terminal, bandara, pelabuhan

---

## Platform & teknologi inti

| Aspek | Target |
|---|---|
| Platform | Windows PC |
| Engine | Unreal Engine 5 |
| Networking | Dedicated Server (server-authoritative) |
| Genre | Open World Multiplayer Life Simulation |
| FPS | 60 FPS stabil di hardware menengah |
| Kapasitas MVP | 20–50 pemain per server |

---

## Pilar gameplay

Pemain dapat menjalani aktivitas sehari-hari tanpa jalan cerita utama yang membatasi. Dunia menjadi tempat pemain menciptakan cerita mereka sendiri.

Contoh aktivitas: jalan, lari, duduk, tidur, makan/minum, mengemudi (mobil/motor/bus), memancing, berkebun, belanja, menabung, beli rumah/kendaraan, buka bisnis, ambil pekerjaan, quest, bersosialisasi.

---

## Dunia (Open World)

Map **bukan** replika Indonesia 1:1. Kota fiksi yang menggabungkan nuansa Jakarta, Bandung, Surabaya, Jogja, dan Bali dalam satu wilayah besar.

Zona yang direncanakan (jangka panjang): Downtown, Perumahan, Kampung, Sawah, Pegunungan, Pantai, Pelabuhan, Bandara, Tol, Mall, Pasar Tradisional, fasilitas kota (RS, polisi, pemadam, sekolah, universitas), SPBU/minimarket/bengkel fiksi, salon, gym, cafe, warkop, warnet, wisata, industri, gudang.

---

## Karakter, NPC, kendaraan

- **Karakter**: customizable modular (pria/wanita, kulit, rambut, wajah, pakaian, aksesoris; tinggi/berat opsional).
- **NPC**: hidup via Behavior Tree + State Machine (bukan LLM di tahap awal). Jadwal harian: kerja, makan, belanja, pulang; weekend ke mall; bereaksi hujan, menyeberang, mengemudi, dll.
- **Kendaraan**: motor, mobil, bus, truk, pickup, sepeda, perahu — desain orisinal terinspirasi Indonesia (bukan salin merek). Fitur: lampu, sein, klakson, bensin, kerusakan, perbaikan.

---

## Ekonomi & progresi

- Pekerjaan beragam (ojol, kurir, mekanik, petani, streamer, pemilik warung, dll.) masing-masing dengan loop ekonomi sendiri.
- Ekonomi **fully server-side**: dompet, bank, transfer, ATM, marketplace, trading; harga dinamis; pajak opsional.
- Inventory, rumah (sewa/beli/jual/upgrade/hias), bisnis (termasuk income offline sesuai desain sistem).

---

## Presentasi dunia

- **Cuaca dinamis**: sunny, cloudy, rain, storm, kabut; siklus siang/malam realtime atau dipercepat.
- **Grafik**: target visual tinggi (Nanite, Lumen, VSM, PBR, post-process) dengan optimasi sejak awal. Gameplay > grafik.
- **Audio**: lingkungan Indonesia (burung, lalu lintas, pasar, hujan, pantai, footstep per permukaan; adzan sesuai konteks implementasi yang tepat).

---

## Prinsip pengembangan

1. Gameplay lebih penting daripada grafik.
2. Fondasi backend harus scalable.
3. Semua sistem penting bersifat server-side.
4. Optimasi dilakukan sejak awal.
5. Asset dapat diganti bertahap tanpa mengubah gameplay.
6. Dunia harus terasa hidup, bukan hanya besar.
7. Setiap update harus menambah nilai bagi pemain.
8. Hindari aset/merek/logo yang melanggar hak cipta; gunakan desain orisinal.
9. Kode modular, terdokumentasi, mudah dikembangkan tim.
10. Target akhir: dunia virtual Indonesia yang terus berkembang bertahun-tahun.

---

## Referensi terkait

- [MVP 0.1](./MVP.md) — ruang lingkup yang dikunci untuk rilis pertama
- [Architecture](./ARCHITECTURE.md) — fondasi teknis
- [Roadmap](./ROADMAP.md) — urutan milestone
- [Tech Decisions](./TECH_DECISIONS.md) — keputusan stack & alasan
