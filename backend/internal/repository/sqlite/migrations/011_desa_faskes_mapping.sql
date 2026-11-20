-- Mapping desa → puskesmas domisili. Satu desa hanya boleh punya satu
-- mapping aktif; perubahan mapping menutup baris lama dan menambah versi baru.
CREATE TABLE IF NOT EXISTS desa_faskes_mapping (
  id TEXT PRIMARY KEY,
  desa_id TEXT NOT NULL REFERENCES wilayah(id),
  faskes_id TEXT NOT NULL REFERENCES faskes(id),
  berlaku_dari TEXT NOT NULL,
  berlaku_sampai TEXT,
  aktif INTEGER NOT NULL DEFAULT 1,
  created_by TEXT REFERENCES users(id),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_desa_faskes_active
  ON desa_faskes_mapping(desa_id)
  WHERE aktif = 1 AND berlaku_sampai IS NULL;

CREATE INDEX IF NOT EXISTS ix_desa_faskes_faskes
  ON desa_faskes_mapping(faskes_id, aktif);

INSERT OR IGNORE INTO desa_faskes_mapping
  (id, desa_id, faskes_id, berlaku_dari, berlaku_sampai, aktif, created_by, created_at, updated_at)
VALUES
  ('map-desa-dalam-taliwang', 'wil-desa-dalam', 'faskes-pkm-taliwang',
   '2026-01-01', NULL, 1, 'usr-admin', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
  ('map-desa-beru-taliwang', 'wil-desa-beru', 'faskes-pkm-taliwang',
   '2026-01-01', NULL, 1, 'usr-admin', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z');
