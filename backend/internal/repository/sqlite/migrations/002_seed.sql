INSERT OR IGNORE INTO users (id, nama, role, faskes_id, bidang_id, aktif, created_at) VALUES
  ('usr-siti',    'Siti',        'pkm_operator', 'faskes-pkm-taliwang', NULL, 1, '2026-05-01T08:00:00Z'),
  ('usr-rina',    'Rina',        'dinas_verifier', NULL, 'bid-kesmas', 1, '2026-05-01T08:00:00Z'),
  ('usr-kabid',   'dr. Ahmad',   'kabid',        NULL, 'bid-kesmas', 1, '2026-05-01T08:00:00Z'),
  ('usr-sekretaris', 'Sekretaris Dinas', 'sekretaris', NULL, 'bid-sekretariat', 1, '2026-05-01T08:00:00Z'),
  ('usr-pimpinan','Kepala Dinas','pimpinan',     NULL, NULL, 1, '2026-05-01T08:00:00Z'),
  ('usr-admin',   'Admin MATA SIDIK', 'admin',   NULL, NULL, 1, '2026-05-01T08:00:00Z'),
  ('usr-pj-gizi', 'PJ Gizi',     'pkm_pj',       NULL, 'bid-gizi', 1, '2026-05-01T08:00:00Z');

INSERT OR IGNORE INTO wilayah (id, kode, nama, level, parent_id) VALUES
  ('wil-kab-swb', '5201', 'Kabupaten Sumbawa Barat', 'kabupaten', NULL),
  ('wil-kec-taliwang', '5201.01', 'Kecamatan Taliwang', 'kecamatan', 'wil-kab-swb'),
  ('wil-kec-jereweh', '5201.02', 'Kecamatan Jereweh', 'kecamatan', 'wil-kab-swb'),
  ('wil-desa-dalam', '5201.01.01', 'Desa Dalam', 'desa', 'wil-kec-taliwang'),
  ('wil-desa-beru', '5201.01.02', 'Desa Beru', 'desa', 'wil-kec-taliwang');

INSERT OR IGNORE INTO faskes (id, kode, nama, wilayah_id, jenis, aktif) VALUES
  ('faskes-pkm-taliwang', 'PKM-TALIWANG', 'PKM Taliwang', 'wil-kec-taliwang', 'puskesmas', 1),
  ('faskes-pkm-jereweh',  'PKM-JEREWEH',  'PKM Jereweh',  'wil-kec-jereweh',  'puskesmas', 1);

INSERT OR IGNORE INTO bidang (id, kode, nama) VALUES
  ('bid-kesmas', 'KESMAS', 'Bidang Kesehatan Masyarakat'),
  ('bid-gizi', 'GIZI', 'Timker Gizi'),
  ('bid-kesga', 'KESGA', 'Timker Kesehatan Keluarga'),
  ('bid-sekretariat', 'SEKRETARIAT', 'Sekretariat Dinas');

INSERT OR IGNORE INTO periode (id, tahun, bulan, kode, status, is_late_allowed) VALUES
  ('per-2026-05', 2026, 5, '2026-05', 'open', 0),
  ('per-2026-06', 2026, 6, '2026-06', 'open', 0);

INSERT OR IGNORE INTO master_release (id, tahun, status, published_at, published_by) VALUES
  ('mr-2026', 2026, 'published', '2026-01-02T00:00:00Z', 'usr-admin');

INSERT OR IGNORE INTO workflow_version (id, kode, version_no, steps_json, status, published_at) VALUES
  ('wf-default-v1', 'default_3_level', 1,
   '[{"kode":"DRAFT_PKM","urutan":1,"role":"pkm_operator"},{"kode":"VERIFY_DINAS","urutan":2,"role":"dinas_verifier"},{"kode":"VALIDATE_KABID","urutan":3,"role":"kabid"},{"kode":"FINAL","urutan":4,"role":"sistem"}]',
   'published', '2026-01-02T00:00:00Z');

INSERT OR IGNORE INTO label
  (id, kode, nama, tipe, parent_id, master_release_id, aggregation_method,
   allow_direct_input, is_rollup_only, pj_bidang_id, iku_owner_id, level, path, aktif, berlaku_dari)
VALUES
  ('lbl-iku-001', 'IKU-2026-001', 'Kesehatan Ibu', 'IKU', NULL, 'mr-2026', 'RATIO_OF_SUMS', 0, 1, 'bid-kesmas', 'usr-kabid', 1, 'IKU-2026-001', 1, '2026-01-01'),
  ('lbl-iku-001-01', 'IKU-2026-001.01', 'Ibu Hamil', 'IKU', 'lbl-iku-001', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'IKU-2026-001/IKU-2026-001.01', 1, '2026-01-01'),
  ('lbl-iku-001-02', 'IKU-2026-001.02', 'Ibu Bersalin', 'IKU', 'lbl-iku-001', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'IKU-2026-001/IKU-2026-001.02', 1, '2026-01-01'),
  ('lbl-iku-001-03', 'IKU-2026-001.03', 'Bayi Baru Lahir', 'IKU', 'lbl-iku-001', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'IKU-2026-001/IKU-2026-001.03', 1, '2026-01-01'),
  ('lbl-spm-2026', 'SPM-2026', 'Standar Pelayanan Minimal 2026', 'SPM', NULL, 'mr-2026', 'RATIO_OF_SUMS', 0, 1, 'bid-kesmas', 'usr-kabid', 1, 'SPM-2026', 1, '2026-01-01'),
  ('lbl-spm-04', 'SPM-2026-04', 'Pelayanan Kesehatan Balita', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-04', 1, '2026-01-01'),
  ('lbl-spm-05', 'SPM-2026-05', 'Pelayanan Kesehatan Usia Pendidikan Dasar', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-05', 1, '2026-01-01'),
  ('lbl-spm-06', 'SPM-2026-06', 'Pelayanan Kesehatan Usia Produktif', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-06', 1, '2026-01-01'),
  ('lbl-spm-07', 'SPM-2026-07', 'Pelayanan Kesehatan Usia Lanjut', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-07', 1, '2026-01-01'),
  ('lbl-spm-08', 'SPM-2026-08', 'Pelayanan Kesehatan Penderita Hipertensi', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-08', 1, '2026-01-01'),
  ('lbl-spm-09', 'SPM-2026-09', 'Pelayanan Kesehatan Penderita Diabetes Melitus', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-09', 1, '2026-01-01'),
  ('lbl-spm-10', 'SPM-2026-10', 'Pelayanan Kesehatan ODGJ Berat', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-10', 1, '2026-01-01'),
  ('lbl-spm-11', 'SPM-2026-11', 'Pelayanan Kesehatan Terduga Tuberkulosis', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-11', 1, '2026-01-01'),
  ('lbl-spm-12', 'SPM-2026-12', 'Pelayanan Kesehatan Risiko Terinfeksi HIV', 'SPM', 'lbl-spm-2026', 'mr-2026', 'RATIO_OF_SUMS', 1, 0, 'bid-kesmas', 'usr-kabid', 2, 'SPM-2026/SPM-2026-12', 1, '2026-01-01');

INSERT OR IGNORE INTO indikator
  (id, kode, nama, satuan, definisi, numerator_label, denominator_label, sifat_input, frekuensi, pj_bidang_id, pj_timker, aktif)
VALUES
  ('ind-001', 'IND-001', 'Jumlah ibu hamil dilayani sesuai standar', 'orang', 'Jumlah ibu hamil yang dilayani sesuai standar BNBA.', 'Jumlah ibu hamil dilayani', 'Jumlah ibu hamil sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-002', 'IND-002', 'Jumlah ibu bersalin dilayani sesuai standar', 'orang', 'Jumlah ibu bersalin yang dilayani sesuai standar BNBA.', 'Jumlah ibu bersalin dilayani', 'Jumlah ibu bersalin sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-003', 'IND-003', 'Jumlah bayi baru lahir dilayani sesuai standar', 'orang', 'Jumlah bayi baru lahir yang dilayani sesuai standar BNBA.', 'Jumlah bayi baru lahir dilayani', 'Jumlah bayi baru lahir sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-004', 'IND-004', 'Pelayanan Kesehatan Balita', '%', 'Cakupan pelayanan kesehatan balita sesuai standar.', 'Balita dilayani', 'Balita sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Gizi', 1),
  ('ind-005', 'IND-005', 'Pelayanan Kesehatan Usia Pendidikan Dasar', '%', 'Cakupan pelayanan kesehatan usia pendidikan dasar.', 'Anak usia pendidikan dasar dilayani', 'Anak usia pendidikan dasar sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-006', 'IND-006', 'Pelayanan Kesehatan Usia Produktif', '%', 'Cakupan pelayanan kesehatan usia produktif.', 'Usia produktif dilayani', 'Usia produktif sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-007', 'IND-007', 'Pelayanan Kesehatan Usia Lanjut', '%', 'Cakupan pelayanan kesehatan usia lanjut.', 'Usia lanjut dilayani', 'Usia lanjut sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-008', 'IND-008', 'Pelayanan Kesehatan Penderita Hipertensi', '%', 'Cakupan pelayanan penderita hipertensi sesuai standar.', 'Penderita hipertensi dilayani', 'Penderita hipertensi sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-009', 'IND-009', 'Pelayanan Kesehatan Penderita Diabetes Melitus', '%', 'Cakupan pelayanan penderita diabetes melitus sesuai standar.', 'Penderita DM dilayani', 'Penderita DM sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-010', 'IND-010', 'Pelayanan Kesehatan ODGJ Berat', '%', 'Cakupan pelayanan orang dengan gangguan jiwa berat.', 'ODGJ berat dilayani', 'ODGJ berat sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-011', 'IND-011', 'Pelayanan Kesehatan Terduga Tuberkulosis', '%', 'Cakupan pelayanan orang terduga tuberkulosis.', 'Terduga TB dilayani', 'Terduga TB sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1),
  ('ind-012', 'IND-012', 'Pelayanan Kesehatan Risiko Terinfeksi HIV', '%', 'Cakupan pelayanan orang dengan risiko terinfeksi HIV.', 'Risiko HIV dilayani', 'Risiko HIV sasaran', 'puskesmas', 'bulanan', 'bid-kesmas', 'Timker Kesga', 1);

INSERT OR IGNORE INTO indikator_label (indikator_id, label_id, bobot, is_primary) VALUES
  ('ind-001', 'lbl-iku-001-01', 1, 1),
  ('ind-002', 'lbl-iku-001-02', 1, 1),
  ('ind-003', 'lbl-iku-001-03', 1, 1),
  ('ind-004', 'lbl-spm-04', 1, 1),
  ('ind-005', 'lbl-spm-05', 1, 1),
  ('ind-006', 'lbl-spm-06', 1, 1),
  ('ind-007', 'lbl-spm-07', 1, 1),
  ('ind-008', 'lbl-spm-08', 1, 1),
  ('ind-009', 'lbl-spm-09', 1, 1),
  ('ind-010', 'lbl-spm-10', 1, 1),
  ('ind-011', 'lbl-spm-11', 1, 1),
  ('ind-012', 'lbl-spm-12', 1, 1);

INSERT OR IGNORE INTO iku_komponen (id, parent_label_id, tipe_komponen, komponen_id, bobot, sort) VALUES
  ('komp-001', 'lbl-iku-001', 'CHILD_IKU', 'lbl-iku-001-01', 1, 1),
  ('komp-002', 'lbl-iku-001', 'CHILD_IKU', 'lbl-iku-001-02', 1, 2),
  ('komp-003', 'lbl-iku-001', 'CHILD_IKU', 'lbl-iku-001-03', 1, 3);

INSERT OR IGNORE INTO target
  (id, label_id, indikator_id, periode_id, scope_type, scope_id, target_tahunan, target_bulanan)
VALUES
  ('tgt-iku-001', 'lbl-iku-001', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-iku-001-01', 'lbl-iku-001-01', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-iku-001-02', 'lbl-iku-001-02', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-iku-001-03', 'lbl-iku-001-03', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-04', 'lbl-spm-04', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-05', 'lbl-spm-05', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-06', 'lbl-spm-06', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-07', 'lbl-spm-07', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-08', 'lbl-spm-08', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-09', 'lbl-spm-09', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-10', 'lbl-spm-10', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-11', 'lbl-spm-11', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100),
  ('tgt-spm-12', 'lbl-spm-12', NULL, 'per-2026-05', 'KABUPATEN', 'wil-kab-swb', 100, 100);

INSERT OR IGNORE INTO dashboard_row (id, kode, label_id, indikator_id, urutan, target, baseline_json, format_json) VALUES
  ('dash-01', '01', 'lbl-iku-001-01', 'ind-001', 1, 100, '{"capaian_2025":100,"capaian_tw1":27.3}', '{"suffix":"%"}'),
  ('dash-02', '02', 'lbl-iku-001-02', 'ind-002', 2, 100, '{"capaian_2025":96,"capaian_tw1":26.2}', '{"suffix":"%"}'),
  ('dash-03', '03', 'lbl-iku-001-03', 'ind-003', 3, 100, '{"capaian_2025":100,"capaian_tw1":28.1}', '{"suffix":"%"}'),
  ('dash-04', '04', 'lbl-spm-04', 'ind-004', 4, 100, '{"capaian_2025":94,"capaian_tw1":22.3}', '{"suffix":"%"}'),
  ('dash-05', '05', 'lbl-spm-05', 'ind-005', 5, 100, '{"capaian_2025":94,"capaian_tw1":75.1}', '{"suffix":"%"}'),
  ('dash-06', '06', 'lbl-spm-06', 'ind-006', 6, 100, '{"capaian_2025":100,"capaian_tw1":17.3}', '{"suffix":"%"}'),
  ('dash-07', '07', 'lbl-spm-07', 'ind-007', 7, 100, '{"capaian_2025":100,"capaian_tw1":27.2}', '{"suffix":"%"}'),
  ('dash-08', '08', 'lbl-spm-08', 'ind-008', 8, 100, '{"capaian_2025":100,"capaian_tw1":30.2}', '{"suffix":"%"}'),
  ('dash-09', '09', 'lbl-spm-09', 'ind-009', 9, 100, '{"capaian_2025":100,"capaian_tw1":24}', '{"suffix":"%"}'),
  ('dash-10', '10', 'lbl-spm-10', 'ind-010', 10, 100, '{"capaian_2025":90,"capaian_tw1":36.6}', '{"suffix":"%"}'),
  ('dash-11', '11', 'lbl-spm-11', 'ind-011', 11, 100, '{"capaian_2025":100,"capaian_tw1":20.8}', '{"suffix":"%"}'),
  ('dash-12', '12', 'lbl-spm-12', 'ind-012', 12, 100, '{"capaian_2025":100,"capaian_tw1":22.7}', '{"suffix":"%"}');
