# Diagram MATA SIDIK

Dua belas diagram di folder ini dirender dengan [stack-diagram](https://github.com/caru-ini/stack-diagram), bukan Mermaid. Setiap diagram punya satu HTML mandiri (font dan ikon sudah di-inline) dan satu PNG untuk dokumen.

| Diagram | HTML sumber | PNG | Dipakai di |
| --- | --- | --- | --- |
| Alur bisnis 6 tahap | `alur-bisnis.html` | `alur-bisnis.png` | README A.2, `docs/01` §5 |
| Siklus status data | `status-data.html` | `status-data.png` | README A.2 |
| Alur peran (swimlane) | `alur-peran.html` | `alur-peran.png` | README B.2 |
| Model data | `model-data.html` | `model-data.png` | README B.3 |
| Alur data IKU end-to-end + trace 3 baris SPM | `alur-data-iku.html` | `alur-data-iku.png` | README A.7, `docs/04` §10 |
| Pipa upload background | `alur-upload.html` | `alur-upload.png` | `docs/04` §3.2 |
| Workflow approval berjenjang | `alur-workflow.html` | `alur-workflow.png` | `docs/04` §4.1 |
| Provenance IKU induk | `alur-provenance.html` | `alur-provenance.png` | `docs/04` §5.3 |
| Trace output SPM ke file | `alur-trace-spm.html` | `alur-trace-spm.png` | `docs/04` §5.5 |
| Arsitektur aplikasi | `arsitektur-aplikasi.html` | `arsitektur-aplikasi.png` | `docs/04` §6.1, README B.8 |
| Skema database (tabel, kolom, relasi) | `schema-db.html` | `schema-db.png` | `docs/schema.md` |
| Alur login 4 tahap | `alur-login.html` | `alur-login.png` | Referensi desain; halaman login memakai HTML inline |

Dokumen memakai `<img>` biasa dengan satu PNG per diagram. Tidak ada varian dark mode.

> **Catatan status:** `schema-db.png` belum diregenerasi setelah penambahan `dashboard_category`, `label.category_id`, `bidang.jenis`, dan `indikator.pj_timker_id`. Generator `gen-schema.py` sudah diperbarui; jalankan `build.sh` saat dependency stack-diagram/Playwright tersedia.

## Regenerasi

Prasyarat:

- Python 3.9+
- Clone [stack-diagram](https://github.com/caru-ini/stack-diagram)
- Node.js + `playwright` (dengan Google Chrome)
- Ikon Lucide diambil otomatis oleh `fetch_icons.py` (butuh jaringan)

```bash
cd assets/diagrams

# contoh path; sesuaikan dengan mesin Anda
export STACK_DIR=/tmp/opencode/stack-diagram/skills/stack-diagram
export PW_NODE_PATH=/tmp/opencode/pw/node_modules

chmod +x build.sh
./build.sh
```

`build.sh` melakukan: fetch ikon → generate draft dari `gen.py` + `gen-end-to-end.py` + `gen-teknis.py` → inline font/ikon dengan `build.py` → render PNG dengan `render.js`.

Folder `drafts/` dan `icons/` hanya cache lokal (sudah masuk `.gitignore`); yang perlu di-commit adalah HTML mandiri, PNG, serta skrip.

Variabel yang bisa diatur:

| Variabel | Default | Arti |
| --- | --- | --- |
| `STACK_DIR` | `/tmp/opencode/stack-diagram/skills/stack-diagram` | Lokasi skill stack-diagram |
| `ICON_DIR` | `assets/diagrams/icons` | Cache ikon Lucide |
| `PW_NODE_PATH` | `/tmp/opencode/pw/node_modules` | Lokasi modul `playwright` |

## Catatan teknis

- Ukuran canvas: alur-bisnis 1340×600, status-data 1240×460, alur-peran 1360×800, model-data 1320×800, alur-data-iku 1860×1810, alur-upload 1280×620, alur-workflow 1360×620, alur-provenance 1400×560, alur-trace-spm 1860×660, arsitektur-aplikasi 1440×760, schema-db 1880×1946, alur-login 1200×220.
- `render.js` memakai `deviceScaleFactor: 2` dan `fullPage: true` untuk diagram penuh, jadi PNG beresolusi 2×. `alur-login` memakai tinggi tetap 220 agar strip 4 tahap tidak menyisakan ruang kosong. Diagram lebar (`alur-data-iku`, `alur-trace-spm`, `schema-db`) dirender dengan lebar 1860–1920 agar teks panel tetap terbaca.
- `backdrop-filter` pada label node dinonaktifkan di generator (`gen.py`, `gen-end-to-end.py`, `gen-teknis.py`, `gen-schema.py`) untuk menghindari artefak duplikasi ikon saat screenshot Chromium.
- Untuk mengubah layout, edit generator sesuai kelompoknya — `gen.py` (alur bisnis, status, peran, model data), `gen-end-to-end.py` (alur data IKU), `gen-teknis.py` (pipa upload, workflow, provenance, trace SPM, arsitektur), `gen-schema.py` (katalog skema database), atau `gen-login.py` (strip login) — lalu jalankan `build.sh`.
- Ikon memakai Lucide (bukan logo vendor) karena peran/entitas di sini bukan brand. Font: Geist + JetBrains Mono (SIL OFL 1.1). stack-diagram: MIT.
