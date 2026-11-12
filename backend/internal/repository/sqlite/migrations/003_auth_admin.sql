ALTER TABLE users ADD COLUMN username TEXT;
ALTER TABLE users ADD COLUMN password_hash TEXT;
ALTER TABLE users ADD COLUMN position_id TEXT;
ALTER TABLE users ADD COLUMN last_login_at TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS ux_users_username ON users(username);

CREATE TABLE IF NOT EXISTS positions (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL UNIQUE,
  nama TEXT NOT NULL,
  role TEXT NOT NULL,
  level INTEGER NOT NULL DEFAULT 1,
  bidang_id TEXT REFERENCES bidang(id),
  faskes_id TEXT REFERENCES faskes(id),
  aktif INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_session (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  refresh_hash TEXT NOT NULL UNIQUE,
  user_agent TEXT,
  created_at TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  revoked_at TEXT
);
CREATE INDEX IF NOT EXISTS ix_auth_session_user ON auth_session(user_id, revoked_at);

ALTER TABLE audit_log ADD COLUMN impersonator_id TEXT;
CREATE INDEX IF NOT EXISTS ix_audit_impersonator ON audit_log(impersonator_id, created_at);
