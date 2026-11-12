ALTER TABLE target ADD COLUMN target_layanan_tahunan REAL;
ALTER TABLE target ADD COLUMN target_layanan_bulanan REAL;

UPDATE target
SET target_layanan_tahunan = 2000,
    target_layanan_bulanan = 2000.0 / 12.0
WHERE id = 'tgt-iku-001-01';
