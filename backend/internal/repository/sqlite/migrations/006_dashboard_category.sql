CREATE TABLE IF NOT EXISTS dashboard_category (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL UNIQUE,
  nama TEXT NOT NULL,
  jenis TEXT NOT NULL DEFAULT 'LABEL' CHECK (jenis IN ('LABEL', 'OWNER')),
  urutan INTEGER NOT NULL DEFAULT 0,
  aktif INTEGER NOT NULL DEFAULT 1,
  is_system INTEGER NOT NULL DEFAULT 0,
  config_json TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

INSERT OR IGNORE INTO dashboard_category
  (id, kode, nama, jenis, urutan, aktif, is_system, created_at, updated_at)
VALUES
  ('cat-iku', 'IKU', 'IKU (Tujuan dan Sasaran)', 'LABEL', 1, 1, 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
  ('cat-ikk', 'IKK', 'IKK', 'LABEL', 2, 1, 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
  ('cat-spm', 'SPM', 'SPM', 'LABEL', 3, 1, 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
  ('cat-ribk', 'RIBK', 'RIBK', 'LABEL', 4, 1, 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
  ('cat-program', 'PROGRAM', 'Program', 'LABEL', 5, 1, 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
  ('cat-kegiatan', 'KEGIATAN', 'Kegiatan', 'LABEL', 6, 1, 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
  ('cat-sub-kegiatan', 'SUB_KEGIATAN', 'Sub Kegiatan', 'LABEL', 7, 1, 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
  ('cat-bidang-timker', 'BIDANG_TIMKER', 'Bidang/Sekretariat/Tim Kerja', 'OWNER', 8, 1, 1, '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z');

ALTER TABLE label ADD COLUMN category_id TEXT REFERENCES dashboard_category(id);

UPDATE label SET category_id = CASE UPPER(tipe)
  WHEN 'IKU' THEN 'cat-iku'
  WHEN 'IKK' THEN 'cat-ikk'
  WHEN 'SPM' THEN 'cat-spm'
  WHEN 'RIBK' THEN 'cat-ribk'
  WHEN 'PROGRAM' THEN 'cat-program'
  WHEN 'KEGIATAN' THEN 'cat-kegiatan'
  WHEN 'SUB_KEGIATAN' THEN 'cat-sub-kegiatan'
  ELSE NULL
END
WHERE category_id IS NULL;

ALTER TABLE bidang ADD COLUMN jenis TEXT NOT NULL DEFAULT 'bidang';
ALTER TABLE bidang ADD COLUMN parent_id TEXT REFERENCES bidang(id);
ALTER TABLE bidang ADD COLUMN aktif INTEGER NOT NULL DEFAULT 1;

UPDATE bidang SET jenis = 'bidang' WHERE kode = 'KESMAS';
UPDATE bidang SET jenis = 'sekretariat' WHERE kode = 'SEKRETARIAT';
UPDATE bidang SET jenis = 'timker' WHERE kode IN ('GIZI', 'KESGA');

ALTER TABLE indikator ADD COLUMN pj_timker_id TEXT REFERENCES bidang(id);

UPDATE indikator
SET pj_timker_id = (
  SELECT b.id FROM bidang b
  WHERE b.jenis = 'timker'
    AND (
      b.nama = indikator.pj_timker
      OR b.kode = UPPER(REPLACE(TRIM(REPLACE(indikator.pj_timker, 'Timker ', '')), ' ', '_'))
    )
  LIMIT 1
)
WHERE pj_timker_id IS NULL AND pj_timker IS NOT NULL AND TRIM(pj_timker) <> '';

CREATE INDEX IF NOT EXISTS ix_label_category ON label(category_id);
CREATE INDEX IF NOT EXISTS ix_category_urutan ON dashboard_category(aktif, urutan);
CREATE INDEX IF NOT EXISTS ix_bidang_jenis ON bidang(jenis, aktif);
CREATE INDEX IF NOT EXISTS ix_indikator_timker ON indikator(pj_timker_id);
