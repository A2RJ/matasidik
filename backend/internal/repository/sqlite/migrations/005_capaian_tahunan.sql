CREATE TABLE IF NOT EXISTS capaian_tahunan (
  id TEXT PRIMARY KEY,
  label_id TEXT REFERENCES label(id),
  indikator_id TEXT REFERENCES indikator(id),
  tahun INTEGER NOT NULL,
  scope_type TEXT NOT NULL DEFAULT 'KABUPATEN',
  scope_id TEXT,
  nilai REAL NOT NULL,
  catatan TEXT,
  updated_by TEXT REFERENCES users(id),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  CHECK (label_id IS NOT NULL OR indikator_id IS NOT NULL)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_capaian_tahunan_key
  ON capaian_tahunan (
    tahun,
    scope_type,
    COALESCE(scope_id, ''),
    COALESCE(label_id, ''),
    COALESCE(indikator_id, '')
  );
CREATE INDEX IF NOT EXISTS ix_capaian_tahunan_tahun
  ON capaian_tahunan (tahun, scope_type, scope_id);

-- Pindahkan baseline tahunan dari JSON dashboard_row ke tabel terpisah.
-- Capaian TW1 tetap berada di baseline_json untuk saat ini.
INSERT OR IGNORE INTO capaian_tahunan
  (id, label_id, indikator_id, tahun, scope_type, scope_id, nilai, catatan,
   updated_by, created_at, updated_at)
SELECT
  'cap-2025-' || dr.kode,
  dr.label_id,
  dr.indikator_id,
  2025,
  'KABUPATEN',
  'wil-kab-swb',
  CAST(json_extract(dr.baseline_json, '$.capaian_2025') AS REAL),
  NULL,
  'usr-admin',
  '2026-01-02T00:00:00Z',
  '2026-01-02T00:00:00Z'
FROM dashboard_row dr
WHERE dr.baseline_json IS NOT NULL
  AND json_extract(dr.baseline_json, '$.capaian_2025') IS NOT NULL;
