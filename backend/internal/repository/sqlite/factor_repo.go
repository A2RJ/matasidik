package sqlite

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

func (s *Store) ListIKUValuesForFactors(ctx context.Context, periodeID, scopeType, scopeID string) ([]domain.IKUValue, error) {
	return s.ListIKUValues(ctx, periodeID, scopeType, scopeID)
}

func (s *Store) GetFaktor(ctx context.Context, ikuValueID string) (*domain.Faktor, error) {
	var out domain.Faktor
	err := s.db.GetContext(ctx, &out, `SELECT f.id, f.iku_value_id, f.periode_id, f.pendukung, f.penghambat,
		f.created_by, f.updated_at, l.kode AS label_kode, l.nama AS label_nama
		FROM faktor f JOIN iku_value v ON v.id = f.iku_value_id JOIN label l ON l.id = v.label_id
		WHERE f.iku_value_id = ?`, ikuValueID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

func (s *Store) UpsertFaktor(ctx context.Context, f *domain.Faktor) error {
	var existingID string
	_ = s.db.GetContext(ctx, &existingID, `SELECT id FROM faktor WHERE iku_value_id = ?`, f.IKUValueID)
	if existingID == "" {
		f.ID = uuid.NewString()
		_, err := s.db.ExecContext(ctx, `INSERT INTO faktor
			(id, iku_value_id, periode_id, pendukung, penghambat, created_by, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			f.ID, f.IKUValueID, f.PeriodeID, f.Pendukung, f.Penghambat, f.CreatedBy, f.UpdatedAt)
		return err
	}
	f.ID = existingID
	_, err := s.db.ExecContext(ctx, `UPDATE faktor SET pendukung = ?, penghambat = ?, created_by = ?, updated_at = ?
		WHERE id = ?`, f.Pendukung, f.Penghambat, f.CreatedBy, f.UpdatedAt, f.ID)
	return err
}

func (s *Store) ListTindakLanjut(ctx context.Context, periodeID string) ([]domain.TindakLanjut, error) {
	q := `SELECT t.id, t.iku_value_id, t.periode_id, t.masalah, t.rekomendasi, t.pj_user_id,
		t.status, t.created_by, t.created_at, t.updated_at,
		u.nama AS pj_nama, l.kode AS label_kode, l.nama AS label_nama
		FROM tindak_lanjut t
		LEFT JOIN users u ON u.id = t.pj_user_id
		LEFT JOIN iku_value v ON v.id = t.iku_value_id
		LEFT JOIN label l ON l.id = v.label_id`
	args := []any{}
	if periodeID != "" {
		q += ` WHERE t.periode_id = ?`
		args = append(args, periodeID)
	}
	q += ` ORDER BY t.created_at DESC`
	var out []domain.TindakLanjut
	err := s.db.SelectContext(ctx, &out, q, args...)
	return out, err
}

func (s *Store) CreateTindakLanjut(ctx context.Context, tl *domain.TindakLanjut) error {
	if tl.ID == "" {
		tl.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO tindak_lanjut
		(id, iku_value_id, periode_id, masalah, rekomendasi, pj_user_id, status, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tl.ID, tl.IKUValueID, tl.PeriodeID, tl.Masalah, tl.Rekomendasi, tl.PJUserID,
		tl.Status, tl.CreatedBy, tl.CreatedAt, tl.UpdatedAt)
	return err
}

func (s *Store) UpdateTindakLanjut(ctx context.Context, tl *domain.TindakLanjut) error {
	_, err := s.db.ExecContext(ctx, `UPDATE tindak_lanjut SET masalah = ?, rekomendasi = ?,
		pj_user_id = ?, status = ?, updated_at = ? WHERE id = ?`,
		tl.Masalah, tl.Rekomendasi, tl.PJUserID, tl.Status, tl.UpdatedAt, tl.ID)
	return err
}

func (s *Store) GetTindakLanjut(ctx context.Context, id string) (*domain.TindakLanjut, error) {
	var out domain.TindakLanjut
	err := s.db.GetContext(ctx, &out, `SELECT t.id, t.iku_value_id, t.periode_id, t.masalah, t.rekomendasi,
		t.pj_user_id, t.status, t.created_by, t.created_at, t.updated_at
		FROM tindak_lanjut t WHERE t.id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("tindak lanjut")
	}
	return &out, err
}
