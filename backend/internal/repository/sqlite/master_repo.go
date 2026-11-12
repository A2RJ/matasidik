package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

func (s *Store) ListWilayah(ctx context.Context) ([]domain.Wilayah, error) {
	var out []domain.Wilayah
	err := s.db.SelectContext(ctx, &out, `SELECT id, kode, nama, level, parent_id FROM wilayah ORDER BY level, nama`)
	return out, err
}

func (s *Store) ListFaskes(ctx context.Context) ([]domain.Faskes, error) {
	var out []domain.Faskes
	err := s.db.SelectContext(ctx, &out, `SELECT id, kode, nama, wilayah_id, jenis, aktif FROM faskes WHERE aktif = 1 ORDER BY nama`)
	return out, err
}

func (s *Store) GetFaskes(ctx context.Context, id string) (*domain.Faskes, error) {
	var out domain.Faskes
	err := s.db.GetContext(ctx, &out, `SELECT id, kode, nama, wilayah_id, jenis, aktif FROM faskes WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("faskes")
	}
	return &out, err
}

func (s *Store) ListBidang(ctx context.Context) ([]domain.Bidang, error) {
	var out []domain.Bidang
	err := s.db.SelectContext(ctx, &out, `SELECT id, kode, nama FROM bidang ORDER BY nama`)
	return out, err
}

func (s *Store) ListPeriods(ctx context.Context) ([]domain.Periode, error) {
	var out []domain.Periode
	err := s.db.SelectContext(ctx, &out, `SELECT id, tahun, bulan, kode, status, is_late_allowed FROM periode ORDER BY tahun DESC, bulan DESC`)
	return out, err
}

func (s *Store) GetPeriodByKode(ctx context.Context, kode string) (*domain.Periode, error) {
	var out domain.Periode
	err := s.db.GetContext(ctx, &out, `SELECT id, tahun, bulan, kode, status, is_late_allowed FROM periode WHERE kode = ?`, kode)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("periode")
	}
	return &out, err
}

func (s *Store) GetPeriodByID(ctx context.Context, id string) (*domain.Periode, error) {
	var out domain.Periode
	err := s.db.GetContext(ctx, &out, `SELECT id, tahun, bulan, kode, status, is_late_allowed FROM periode WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("periode")
	}
	return &out, err
}

func (s *Store) DefaultPeriodID(ctx context.Context) (string, error) {
	var id string
	err := s.db.GetContext(ctx, &id, `SELECT periode_id FROM submission
		WHERE status = 'final' ORDER BY updated_at DESC LIMIT 1`)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}
	err = s.db.GetContext(ctx, &id, `SELECT id FROM periode ORDER BY tahun DESC, bulan DESC LIMIT 1`)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return id, err
}

func (s *Store) GetPublishedRelease(ctx context.Context, tahun int) (*domain.MasterRelease, error) {
	var out domain.MasterRelease
	err := s.db.GetContext(ctx, &out, `SELECT id, tahun, status, published_at, published_by
		FROM master_release WHERE tahun = ? AND status = 'published' ORDER BY published_at DESC LIMIT 1`, tahun)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("master release")
	}
	return &out, err
}

func (s *Store) GetMasterRelease(ctx context.Context, id string) (*domain.MasterRelease, error) {
	var out domain.MasterRelease
	err := s.db.GetContext(ctx, &out, `SELECT id, tahun, status, published_at, published_by FROM master_release WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("master release")
	}
	return &out, err
}

func (s *Store) GetWorkflowVersion(ctx context.Context, id string) (*domain.WorkflowVersion, error) {
	var out domain.WorkflowVersion
	err := s.db.GetContext(ctx, &out, `SELECT id, kode, version_no, steps_json, status, published_at FROM workflow_version WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("workflow version")
	}
	return &out, err
}

func (s *Store) ListLabels(ctx context.Context) ([]domain.Label, error) {
	var out []domain.Label
	err := s.db.SelectContext(ctx, &out, `SELECT id, kode, nama, tipe, parent_id, master_release_id,
		aggregation_method, allow_direct_input, is_rollup_only, pj_bidang_id, iku_owner_id,
		level, path, aktif, berlaku_dari, berlaku_sampai
		FROM label WHERE aktif = 1 ORDER BY level, kode`)
	return out, err
}

func (s *Store) GetLabelByID(ctx context.Context, id string) (*domain.Label, error) {
	var out domain.Label
	err := s.db.GetContext(ctx, &out, `SELECT id, kode, nama, tipe, parent_id, master_release_id,
		aggregation_method, allow_direct_input, is_rollup_only, pj_bidang_id, iku_owner_id,
		level, path, aktif, berlaku_dari, berlaku_sampai
		FROM label WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("label")
	}
	return &out, err
}

func (s *Store) GetLabelByKode(ctx context.Context, kode, releaseID string) (*domain.Label, error) {
	var out domain.Label
	err := s.db.GetContext(ctx, &out, `SELECT id, kode, nama, tipe, parent_id, master_release_id,
		aggregation_method, allow_direct_input, is_rollup_only, pj_bidang_id, iku_owner_id,
		level, path, aktif, berlaku_dari, berlaku_sampai
		FROM label WHERE kode = ? AND master_release_id = ? AND aktif = 1`, kode, releaseID)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("label")
	}
	return &out, err
}

func (s *Store) ListIndicators(ctx context.Context) ([]domain.Indikator, error) {
	var out []domain.Indikator
	err := s.db.SelectContext(ctx, &out, `SELECT id, kode, nama, satuan, definisi, numerator_label,
		denominator_label, sifat_input, frekuensi, pj_bidang_id, pj_timker, aktif
		FROM indikator WHERE aktif = 1 ORDER BY kode`)
	return out, err
}

func (s *Store) GetIndicatorByID(ctx context.Context, id string) (*domain.Indikator, error) {
	var out domain.Indikator
	err := s.db.GetContext(ctx, &out, `SELECT id, kode, nama, satuan, definisi, numerator_label,
		denominator_label, sifat_input, frekuensi, pj_bidang_id, pj_timker, aktif
		FROM indikator WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("indikator")
	}
	return &out, err
}

func (s *Store) GetIndicatorByKode(ctx context.Context, kode string) (*domain.Indikator, error) {
	var out domain.Indikator
	err := s.db.GetContext(ctx, &out, `SELECT id, kode, nama, satuan, definisi, numerator_label,
		denominator_label, sifat_input, frekuensi, pj_bidang_id, pj_timker, aktif
		FROM indikator WHERE kode = ? AND aktif = 1`, kode)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("indikator")
	}
	return &out, err
}

func (s *Store) ListIndicatorLabels(ctx context.Context, indikatorID string) ([]domain.IndikatorLabel, error) {
	var out []domain.IndikatorLabel
	err := s.db.SelectContext(ctx, &out, `SELECT indikator_id, label_id, bobot, is_primary
		FROM indikator_label WHERE indikator_id = ?`, indikatorID)
	return out, err
}

func (s *Store) ListIndicatorLabelsForLabel(ctx context.Context, labelID string) ([]domain.IndikatorLabel, error) {
	var out []domain.IndikatorLabel
	err := s.db.SelectContext(ctx, &out, `SELECT indikator_id, label_id, bobot, is_primary
		FROM indikator_label WHERE label_id = ?`, labelID)
	return out, err
}

func (s *Store) ListIKUComponents(ctx context.Context, parentLabelID string) ([]domain.IKUComponent, error) {
	var out []domain.IKUComponent
	err := s.db.SelectContext(ctx, &out, `SELECT id, parent_label_id, tipe_komponen, komponen_id, bobot, sort
		FROM iku_komponen WHERE parent_label_id = ? ORDER BY sort`, parentLabelID)
	return out, err
}

func (s *Store) ListTargets(ctx context.Context, periodeID string) ([]domain.Target, error) {
	var out []domain.Target
	q := `SELECT id, label_id, indikator_id, periode_id, scope_type,
		scope_id, target_tahunan, target_bulanan FROM target`
	args := []any{}
	if periodeID != "" {
		q += ` WHERE periode_id = ?`
		args = append(args, periodeID)
	}
	q += ` ORDER BY id`
	err := s.db.SelectContext(ctx, &out, q, args...)
	return out, err
}

func (s *Store) ListDashboardRows(ctx context.Context) ([]domain.DashboardRow, error) {
	var out []domain.DashboardRow
	err := s.db.SelectContext(ctx, &out, `SELECT id, kode, label_id, indikator_id, urutan, target,
		baseline_json, format_json FROM dashboard_row ORDER BY urutan`)
	return out, err
}

func (s *Store) GetDashboardRowByKode(ctx context.Context, kode string) (*domain.DashboardRow, error) {
	var out domain.DashboardRow
	err := s.db.GetContext(ctx, &out, `SELECT id, kode, label_id, indikator_id, urutan, target,
		baseline_json, format_json FROM dashboard_row WHERE kode = ?`, kode)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("baris dashboard")
	}
	return &out, err
}

func (s *Store) CreateIndicator(ctx context.Context, indikator *domain.Indikator, labelIDs []string) error {
	if indikator.ID == "" {
		indikator.ID = uuid.NewString()
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `INSERT INTO indikator
		(id, kode, nama, satuan, definisi, numerator_label, denominator_label, sifat_input, frekuensi, pj_bidang_id, pj_timker, aktif)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		indikator.ID, indikator.Kode, indikator.Nama, indikator.Satuan, indikator.Definisi,
		indikator.NumeratorLabel, indikator.DenominatorLabel, indikator.SifatInput, indikator.Frekuensi,
		indikator.PJBidangID, indikator.PJTimker, indikator.Aktif); err != nil {
		return fmt.Errorf("simpan indikator: %w", err)
	}
	for _, labelID := range labelIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO indikator_label(indikator_id, label_id, bobot, is_primary)
			VALUES (?, ?, 1, 1)`, indikator.ID, labelID); err != nil {
			return fmt.Errorf("hubungkan label: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Store) CreateLabel(ctx context.Context, label *domain.Label) error {
	if label.ID == "" {
		label.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO label
		(id, kode, nama, tipe, parent_id, master_release_id, aggregation_method, allow_direct_input,
		 is_rollup_only, pj_bidang_id, iku_owner_id, level, path, aktif, berlaku_dari, berlaku_sampai)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		label.ID, label.Kode, label.Nama, label.Tipe, label.ParentID, label.MasterReleaseID,
		label.AggregationMethod, label.AllowDirectInput, label.IsRollupOnly, label.PJBidangID,
		label.IKUOwnerID, label.Level, label.Path, label.Aktif, label.BerlakuDari, label.BerlakuSampai)
	return err
}

func (s *Store) UpsertTarget(ctx context.Context, target *domain.Target) error {
	var existingID string
	err := s.db.GetContext(ctx, &existingID, `SELECT id FROM target
		WHERE COALESCE(label_id,'') = COALESCE(?,'')
		  AND COALESCE(indikator_id,'') = COALESCE(?,'')
		  AND periode_id = ? AND scope_type = ? AND COALESCE(scope_id,'') = COALESCE(?,'')`,
		target.LabelID, target.IndikatorID, target.PeriodeID, target.ScopeType, target.ScopeID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingID != "" {
		target.ID = existingID
		_, err = s.db.ExecContext(ctx, `UPDATE target SET target_tahunan = ?, target_bulanan = ? WHERE id = ?`,
			target.TargetTahunan, target.TargetBulanan, target.ID)
		return err
	}
	if target.ID == "" {
		target.ID = uuid.NewString()
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO target
		(id, label_id, indikator_id, periode_id, scope_type, scope_id, target_tahunan, target_bulanan)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		target.ID, target.LabelID, target.IndikatorID, target.PeriodeID, target.ScopeType,
		target.ScopeID, target.TargetTahunan, target.TargetBulanan)
	return err
}
