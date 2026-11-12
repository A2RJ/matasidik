package sqlite

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

const capaianSelect = `SELECT c.id, c.label_id, c.indikator_id, c.tahun, c.scope_type, c.scope_id,
	c.nilai, c.catatan, c.updated_by, c.created_at, c.updated_at,
	l.kode AS label_kode, l.nama AS label_nama,
	i.kode AS indikator_kode, i.nama AS indikator_nama,
	u.nama AS updated_by_nama
	FROM capaian_tahunan c
	LEFT JOIN label l ON l.id = c.label_id
	LEFT JOIN indikator i ON i.id = c.indikator_id
	LEFT JOIN users u ON u.id = c.updated_by`

func (s *Store) ListCapaianTahunan(ctx context.Context, tahun int, scopeType, scopeID string) ([]domain.CapaianTahunan, error) {
	var out []domain.CapaianTahunan
	err := s.db.SelectContext(ctx, &out, capaianSelect+`
		WHERE c.tahun = ? AND c.scope_type = ? AND COALESCE(c.scope_id,'') = COALESCE(?,'')
		ORDER BY COALESCE(l.kode, i.kode), c.created_at`, tahun, scopeType, scopeID)
	return out, err
}

func (s *Store) GetCapaianTahunan(ctx context.Context, id string) (*domain.CapaianTahunan, error) {
	var out domain.CapaianTahunan
	err := s.db.GetContext(ctx, &out, capaianSelect+` WHERE c.id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("capaian tahunan")
	}
	return &out, err
}

func (s *Store) UpsertCapaianTahunan(ctx context.Context, c *domain.CapaianTahunan) error {
	if c.ID == "" {
		var existingID string
		err := s.db.GetContext(ctx, &existingID, `SELECT id FROM capaian_tahunan
			WHERE tahun = ? AND scope_type = ? AND COALESCE(scope_id,'') = COALESCE(?,'')
			  AND COALESCE(label_id,'') = COALESCE(?,'')
			  AND COALESCE(indikator_id,'') = COALESCE(?,'')`,
			c.Tahun, c.ScopeType, c.ScopeID, c.LabelID, c.IndikatorID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		c.ID = existingID
	}
	if c.ID != "" {
		result, err := s.db.ExecContext(ctx, `UPDATE capaian_tahunan SET
			label_id = ?, indikator_id = ?, tahun = ?, scope_type = ?, scope_id = ?,
			nilai = ?, catatan = ?, updated_by = ?, updated_at = ? WHERE id = ?`,
			c.LabelID, c.IndikatorID, c.Tahun, c.ScopeType, c.ScopeID,
			c.Nilai, c.Catatan, c.UpdatedBy, c.UpdatedAt, c.ID)
		if err != nil {
			return err
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return domain.NotFound("capaian tahunan")
		}
		return nil
	}
	c.ID = uuid.NewString()
	if c.CreatedAt == "" {
		c.CreatedAt = nowRFC3339()
	}
	if c.UpdatedAt == "" {
		c.UpdatedAt = c.CreatedAt
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO capaian_tahunan
		(id, label_id, indikator_id, tahun, scope_type, scope_id, nilai, catatan,
		 updated_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.LabelID, c.IndikatorID, c.Tahun, c.ScopeType, c.ScopeID,
		c.Nilai, c.Catatan, c.UpdatedBy, c.CreatedAt, c.UpdatedAt)
	return err
}

func (s *Store) DeleteCapaianTahunan(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM capaian_tahunan WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.NotFound("capaian tahunan")
	}
	return nil
}
