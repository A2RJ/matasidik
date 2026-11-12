INSERT OR IGNORE INTO positions (id, kode, nama, role, level, bidang_id, faskes_id, aktif, created_at) VALUES
  ('pos-admin', 'ADMIN', 'Administrator Sistem', 'admin', 1, NULL, NULL, 1, '2026-01-02T00:00:00Z'),
  ('pos-pimpinan', 'PIMPINAN', 'Kepala Dinas', 'pimpinan', 1, NULL, NULL, 1, '2026-01-02T00:00:00Z'),
  ('pos-verifikator', 'VERIFIKATOR', 'Verifikator Dinas', 'dinas_verifier', 2, 'bid-kesmas', NULL, 1, '2026-01-02T00:00:00Z'),
  ('pos-kabid-kesmas', 'KABID-KESMAS', 'Kepala Bidang Kesehatan Masyarakat', 'kabid', 3, 'bid-kesmas', NULL, 1, '2026-01-02T00:00:00Z'),
  ('pos-sekretaris', 'SEKRETARIS', 'Sekretaris Dinas', 'sekretaris', 3, 'bid-sekretariat', NULL, 1, '2026-01-02T00:00:00Z'),
  ('pos-operator-taliwang', 'OPERATOR-TALIWANG', 'Operator PKM Taliwang', 'pkm_operator', 1, NULL, 'faskes-pkm-taliwang', 1, '2026-01-02T00:00:00Z'),
  ('pos-pj-gizi', 'PJ-GIZI', 'PJ Program Gizi', 'pkm_pj', 1, 'bid-gizi', NULL, 1, '2026-01-02T00:00:00Z');

-- Password pilot: password123 (bcrypt cost 12). Ganti sebelum produksi.
UPDATE users SET username = 'siti', password_hash = '$2a$12$ZyNszKJDkWth/Q5bv8xBJ..ZV6LCXbIrsaS4TkRs3TWa7iR8PKCZy', position_id = 'pos-operator-taliwang' WHERE id = 'usr-siti';
UPDATE users SET username = 'rina', password_hash = '$2a$12$ZyNszKJDkWth/Q5bv8xBJ..ZV6LCXbIrsaS4TkRs3TWa7iR8PKCZy', position_id = 'pos-verifikator' WHERE id = 'usr-rina';
UPDATE users SET username = 'kabid', password_hash = '$2a$12$ZyNszKJDkWth/Q5bv8xBJ..ZV6LCXbIrsaS4TkRs3TWa7iR8PKCZy', position_id = 'pos-kabid-kesmas' WHERE id = 'usr-kabid';
UPDATE users SET username = 'sekretaris', password_hash = '$2a$12$ZyNszKJDkWth/Q5bv8xBJ..ZV6LCXbIrsaS4TkRs3TWa7iR8PKCZy', position_id = 'pos-sekretaris' WHERE id = 'usr-sekretaris';
UPDATE users SET username = 'pimpinan', password_hash = '$2a$12$ZyNszKJDkWth/Q5bv8xBJ..ZV6LCXbIrsaS4TkRs3TWa7iR8PKCZy', position_id = 'pos-pimpinan' WHERE id = 'usr-pimpinan';
UPDATE users SET username = 'admin', password_hash = '$2a$12$ZyNszKJDkWth/Q5bv8xBJ..ZV6LCXbIrsaS4TkRs3TWa7iR8PKCZy', position_id = 'pos-admin' WHERE id = 'usr-admin';
UPDATE users SET username = 'pjgizi', password_hash = '$2a$12$ZyNszKJDkWth/Q5bv8xBJ..ZV6LCXbIrsaS4TkRs3TWa7iR8PKCZy', position_id = 'pos-pj-gizi' WHERE id = 'usr-pj-gizi';
