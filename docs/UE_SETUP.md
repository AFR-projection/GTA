# Unreal Engine 5 — Setup (Windows)

**Status mesin (2026-07-19):** Unreal Editor **belum terdeteksi** di PC ini.  
Tanpa UE5, slice multiplayer visual belum bisa dibuat. Backend tetap jalan di Neon.

---

## Tugas anak buah (wajib sebelum “gas unreal” coding)

### 1) Install Epic Games Launcher
https://store.epicgames.com/en-US/download

### 2) Install Unreal Engine **5.4 atau 5.5** (salah satu, jangan campur)
Di Launcher → Unreal Engine → Library → Install Engine.

Komponen saran:
- Engine
- Editor
- Target Platforms: **Windows**
- (Opsional) Starter Content — boleh, nanti bisa dibuang

Disk: sediakan **~100GB+** bebas.

### 3) Verifikasi
Setelah install, path biasanya:

`C:\Program Files\Epic Games\UE_5.X\Engine\Binaries\Win64\UnrealEditor.exe`

Kabari lead: **“UE ready 5.X”** + versi persis.

---

## Buat project game di repo

Setelah UE terpasang:

1. Buka Unreal Editor  
2. **Games → Third Person** (template bagus buat movement awal)  
3. Project Location:  
   `C:\Users\User\Documents\aldopr\game`  
4. Project Name: `ILO`  
5. Centang **C++** jika sudah punya Visual Studio 2022 + workload “Game development with C++”  
   - Kalau belum VS: boleh **Blueprint** dulu, C++ menyusul  
6. Create

Hasil yang diharapkan di git:

```
game/
  ILO.uproject
  Content/
  Config/
  Source/          (jika C++)
  README.md
```

**Jangan commit:** `Binaries/`, `DerivedDataCache/`, `Intermediate/`, `Saved/` (sudah di `.gitignore`).

### 7) Dedicated server (setelah project ada)
Ikuti [UE_INTEGRATION.md](./UE_INTEGRATION.md) — target pertama:

1. 1 dedicated server + 2 client  
2. Movement sync  
3. Global chat text  
4. Baru HTTP login ke Go API  

---

## Visual Studio (untuk C++)

Kalau pilih C++ project:

1. Install Visual Studio 2022 Community  
2. Workload: **Game development with C++**  
3. Pastikan Windows SDK terpasang  

Blueprint-only boleh untuk minggu 1 multiplayer smoke — tapi jangka panjang kita butuh C++ untuk networking/plugins bersih.

---

## Jangan lakukan

- Jangan bikin map raksasa sebelum 2 client sync movement  
- Jangan customizer karakter ultra sebelum login→spawn→save  
- Jangan salin aset merek nyata  

Lihat [MVP.md](./MVP.md) dan [PRD.md](./PRD.md).
