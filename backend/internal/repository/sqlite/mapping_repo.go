package sqlite

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

const desaFaskesSelect = `SELECT m.id, m.desa_id, m.faskes_id, m.berlaku_dari,
	m.berlaku_sampai, m.aktif, m.created_by, m.created_at, m.updated_at,
	d.kode AS desa_kode, d.nama AS desa_nama,
	f.kode AS faskes_kode, f.nama AS faskes_nama
	FROM desa_faskes_mapping m
	JOIN wilayah d ON d.id = m.desa_id
	JOIN faskes f ON f.id = m.faskes_id`

func (s *Store) ListDesaFaskesMappings(ctx context.Context) ([]domain.DesaFaskesMapping, error) {
	var out []domain.DesaFaskesMapping
	err := s.db.SelectContext(ctx, &out, desaFaskesSelect+`
		ORDER BY m.aktif DESC, d.kode, m.berlaku_dari DESC`)
	return out, err
}

func (s *Store) GetActiveFaskesForDesa(ctx context.Context, desaID string) (*domain.DesaFaskesMapping, error) {
	var out domain.DesaFaskesMapping
	err := s.db.GetContext(ctx, &out, desaFaskesSelect+`
		WHERE m.desa_id = ? AND m.aktif = 1
		  AND m.berlaku_dari <= date('now')
		  AND (m.berlaku_sampai IS NULL OR m.berlaku_sampai > date('now'))
		ORDER BY m.berlaku_dari DESC LIMIT 1`, desaID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *Store) ListActiveDesaIDsForFaskes(ctx context.Context, faskesID string) ([]string, error) {
	var out []string
	err := s.db.SelectContext(ctx, &out, `SELECT desa_id FROM desa_faskes_mapping
		WHERE faskes_id = ? AND aktif = 1
		  AND berlaku_dari <= date('now')
		  AND (berlaku_sampai IS NULL OR berlaku_sampai > date('now'))
		ORDER BY desa_id`, faskesID)
	return out, err
}

// UpsertDesaFaskesMapping menutup mapping aktif lama untuk desa yang sama,
// lalu menyimpan versi baru. Pemanggil mengisi BerlakuDari/UpdatedAt.
func (s *Store) UpsertDesaFaskesMapping(ctx context.Context, m *domain.DesaFaskesMapping) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `UPDATE desa_faskes_mapping
		SET aktif = 0, berlaku_sampai = ?, updated_at = ?
		WHERE desa_id = ? AND aktif = 1 AND berlaku_sampai IS NULL`,
		m.BerlakuDari, m.UpdatedAt, m.DesaID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO desa_faskes_mapping
		(id, desa_id, faskes_id, berlaku_dari, berlaku_sampai, aktif, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, NULL, 1, ?, ?, ?)`,
		m.ID, m.DesaID, m.FaskesID, m.BerlakuDari, m.CreatedBy, m.CreatedAt, m.UpdatedAt); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) DeactivateDesaFaskesMapping(ctx context.Context, id, at string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE desa_faskes_mapping
		SET aktif = 0, berlaku_sampai = ?, updated_at = ?
		WHERE id = ? AND aktif = 1`, at, at, id)
	return err
}
