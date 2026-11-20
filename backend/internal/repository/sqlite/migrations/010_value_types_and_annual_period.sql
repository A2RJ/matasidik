-- Nilai indikator tidak selalu rasio. Tipe eksplisit mencegah tebakan dari
-- satuan teks dan memungkinkan count/scalar dihitung tanpa denominator.
ALTER TABLE indikator ADD COLUMN tipe_nilai TEXT NOT NULL DEFAULT 'RATIO'
  CHECK (tipe_nilai IN ('RATIO', 'COUNT', 'SCALAR'));

-- Periode tahunan memakai bulan = 0. SQLite tidak bisa mengubah CHECK lama,
-- jadi tabel periode dibangun ulang. Migrasi ini dijalankan dengan foreign
-- keys dimatikan oleh store.migrate() agar relasi submission/target tetap aman.
CREATE TABLE periode_new (
  id TEXT PRIMARY KEY,
  tahun INTEGER NOT NULL,
  bulan INTEGER NOT NULL CHECK (bulan BETWEEN 0 AND 12),
  kode TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'locked')),
  is_late_allowed INTEGER NOT NULL DEFAULT 0,
  UNIQUE (tahun, bulan)
);

INSERT INTO periode_new (id, tahun, bulan, kode, status, is_late_allowed)
SELECT id, tahun, bulan, kode, status, is_late_allowed FROM periode;

DROP TABLE periode;
ALTER TABLE periode_new RENAME TO periode;

INSERT OR IGNORE INTO periode (id, tahun, bulan, kode, status, is_late_allowed)
VALUES ('per-2026-annual', 2026, 0, '2026', 'open', 0);
