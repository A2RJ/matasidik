package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

const userSelect = `SELECT u.id, u.nama, COALESCE(u.username,'') AS username,
	COALESCE(p.role, u.role) AS role, u.faskes_id, u.bidang_id, u.position_id,
	u.last_login_at, u.aktif, u.created_at, p.nama AS position_nama
	FROM users u LEFT JOIN positions p ON p.id = u.position_id AND p.aktif = 1`

func (s *Store) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := s.db.GetContext(ctx, &u, userSelect+` WHERE u.id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("user")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	err := s.db.SelectContext(ctx, &users, userSelect+` ORDER BY u.nama`)
	return users, err
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*domain.UserAuth, error) {
	var u domain.UserAuth
	err := s.db.GetContext(ctx, &u, `SELECT u.id, u.nama, COALESCE(u.username,'') AS username,
		COALESCE(p.role, u.role) AS role, u.faskes_id, u.bidang_id, u.position_id,
		u.last_login_at, u.aktif, u.created_at, p.nama AS position_nama,
		COALESCE(u.password_hash,'') AS password_hash
		FROM users u LEFT JOIN positions p ON p.id = u.position_id AND p.aktif = 1
		WHERE u.username = ? COLLATE NOCASE`, strings.TrimSpace(username))
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("user")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) ListPositions(ctx context.Context) ([]domain.Position, error) {
	var out []domain.Position
	err := s.db.SelectContext(ctx, &out, `SELECT id, kode, nama, role, level, bidang_id, faskes_id, aktif, created_at
		FROM positions ORDER BY level, nama`)
	return out, err
}

func (s *Store) CreatePosition(ctx context.Context, p *domain.Position) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	if p.CreatedAt == "" {
		p.CreatedAt = nowRFC3339()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO positions
		(id, kode, nama, role, level, bidang_id, faskes_id, aktif, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Kode, p.Nama, p.Role, p.Level, p.BidangID, p.FaskesID, p.Aktif, p.CreatedAt)
	return err
}

func (s *Store) UpdatePosition(ctx context.Context, p *domain.Position) error {
	_, err := s.db.ExecContext(ctx, `UPDATE positions SET kode = ?, nama = ?, role = ?, level = ?,
		bidang_id = ?, faskes_id = ?, aktif = ? WHERE id = ?`,
		p.Kode, p.Nama, p.Role, p.Level, p.BidangID, p.FaskesID, p.Aktif, p.ID)
	return err
}

func (s *Store) CreateUser(ctx context.Context, u *domain.User, passwordHash string) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.CreatedAt == "" {
		u.CreatedAt = nowRFC3339()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO users
		(id, nama, username, password_hash, role, faskes_id, bidang_id, position_id, aktif, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Nama, u.Username, passwordHash, u.Role, u.FaskesID, u.BidangID, u.PositionID, u.Aktif, u.CreatedAt)
	if err != nil {
		return fmt.Errorf("simpan user: %w", err)
	}
	return nil
}

func (s *Store) UpdateUser(ctx context.Context, u *domain.User) error {
	_, err := s.db.ExecContext(ctx, `UPDATE users SET nama = ?, username = ?, faskes_id = ?,
		bidang_id = ?, position_id = ?, aktif = ? WHERE id = ?`,
		u.Nama, u.Username, u.FaskesID, u.BidangID, u.PositionID, u.Aktif, u.ID)
	return err
}

func (s *Store) UpdateUserName(ctx context.Context, id, nama string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE users SET nama = ? WHERE id = ?`, nama, id)
	return err
}

func (s *Store) UpdateUserPassword(ctx context.Context, id, passwordHash string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE users SET password_hash = ? WHERE id = ?`, passwordHash, id)
	return err
}

func (s *Store) CreateFaskes(ctx context.Context, f *domain.Faskes) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO faskes (id, kode, nama, wilayah_id, jenis, aktif)
		VALUES (?, ?, ?, ?, ?, ?)`, f.ID, f.Kode, f.Nama, f.WilayahID, f.Jenis, f.Aktif)
	return err
}

func (s *Store) UpdateFaskes(ctx context.Context, f *domain.Faskes) error {
	_, err := s.db.ExecContext(ctx, `UPDATE faskes SET kode = ?, nama = ?, wilayah_id = ?, jenis = ?, aktif = ?
		WHERE id = ?`, f.Kode, f.Nama, f.WilayahID, f.Jenis, f.Aktif, f.ID)
	return err
}

func (s *Store) CreateSession(ctx context.Context, sess *domain.AuthSession) error {
	if sess.ID == "" {
		sess.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO auth_session
		(id, user_id, refresh_hash, user_agent, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		sess.ID, sess.UserID, sess.RefreshHash, sess.UserAgent, sess.CreatedAt, sess.ExpiresAt)
	return err
}

func (s *Store) GetSessionByHash(ctx context.Context, hash string) (*domain.AuthSession, error) {
	var out domain.AuthSession
	err := s.db.GetContext(ctx, &out, `SELECT id, user_id, refresh_hash, user_agent, created_at,
		expires_at, revoked_at FROM auth_session WHERE refresh_hash = ?`, hash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

func (s *Store) RevokeSession(ctx context.Context, id, revokedAt string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE auth_session SET revoked_at = ? WHERE id = ?`, revokedAt, id)
	return err
}

func (s *Store) RevokeSessionsForUser(ctx context.Context, userID, revokedAt string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE auth_session SET revoked_at = ?
		WHERE user_id = ? AND revoked_at IS NULL`, revokedAt, userID)
	return err
}

func (s *Store) UpdateLastLogin(ctx context.Context, userID, at string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE users SET last_login_at = ? WHERE id = ?`, at, userID)
	return err
}
