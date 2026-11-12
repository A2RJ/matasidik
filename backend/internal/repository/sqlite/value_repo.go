package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

func (s *Store) ListFinalIndicatorValues(ctx context.Context, periodeID string) ([]domain.IndicatorValue, error) {
	return s.ListFinalIndicatorValuesForScope(ctx, periodeID, "", "")
}

func (s *Store) ListFinalIndicatorValuesForScope(ctx context.Context, periodeID, scopeType, scopeID string) ([]domain.IndicatorValue, error) {
	q := `SELECT iv.id, iv.submission_row_id, iv.indikator_id, iv.periode_id, iv.faskes_id,
		iv.scope_type, iv.scope_id, iv.numerator, iv.denominator, iv.value, iv.status,
		iv.data_version, iv.created_at
		FROM indicator_value iv
		JOIN submission_row sr ON sr.id = iv.submission_row_id
		JOIN submission s ON s.id = sr.submission_id
		WHERE s.status = 'final' AND iv.periode_id = ?`
	args := []any{periodeID}
	if scopeType == "FASKES" && scopeID != "" {
		q += ` AND iv.faskes_id = ?`
		args = append(args, scopeID)
	}
	q += ` ORDER BY iv.created_at`
	var out []domain.IndicatorValue
	err := s.db.SelectContext(ctx, &out, q, args...)
	return out, err
}

func (s *Store) GetIKUValue(ctx context.Context, labelID, periodeID, scopeType, scopeID string) (*domain.IKUValue, error) {
	var out domain.IKUValue
	err := s.db.GetContext(ctx, &out, `SELECT v.id, v.label_id, v.periode_id, v.scope_type, v.scope_id,
		v.numerator, v.denominator, v.value, v.status, v.formula_version, v.calculation_run_id,
		v.source_hash, v.data_version, v.updated_at,
		l.kode AS label_kode, l.nama AS label_nama, l.tipe AS label_tipe
		FROM iku_value v JOIN label l ON l.id = v.label_id
		WHERE v.label_id = ? AND v.periode_id = ? AND v.scope_type = ? AND COALESCE(v.scope_id,'') = COALESCE(?,'')`,
		labelID, periodeID, scopeType, scopeID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

func (s *Store) ListIKUValues(ctx context.Context, periodeID, scopeType, scopeID string) ([]domain.IKUValue, error) {
	var out []domain.IKUValue
	err := s.db.SelectContext(ctx, &out, `SELECT v.id, v.label_id, v.periode_id, v.scope_type, v.scope_id,
		v.numerator, v.denominator, v.value, v.status, v.formula_version, v.calculation_run_id,
		v.source_hash, v.data_version, v.updated_at,
		l.kode AS label_kode, l.nama AS label_nama, l.tipe AS label_tipe
		FROM iku_value v JOIN label l ON l.id = v.label_id
		WHERE v.periode_id = ? AND v.scope_type = ? AND COALESCE(v.scope_id,'') = COALESCE(?,'')
		ORDER BY l.level, l.kode`, periodeID, scopeType, scopeID)
	return out, err
}

func (s *Store) ListIKUValueDetails(ctx context.Context, ikuValueID string) ([]domain.IKUValueDetail, error) {
	var out []domain.IKUValueDetail
	err := s.db.SelectContext(ctx, &out, `SELECT d.id, d.iku_value_id, d.component_type, d.component_id,
		d.component_label_id, d.component_value, d.numerator, d.denominator, d.weight, d.contribution,
		l.kode AS label_kode, l.nama AS label_nama, i.kode AS indikator_kode, i.nama AS indikator_nama
		FROM iku_value_detail d
		LEFT JOIN label l ON l.id = d.component_label_id
		LEFT JOIN indikator i ON i.id = d.component_id
		WHERE d.iku_value_id = ? ORDER BY d.id`, ikuValueID)
	return out, err
}

func (s *Store) UpsertIKUValues(ctx context.Context, writes []domain.IKUValueWrite) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, write := range writes {
		v := write.Value
		var existingID string
		_ = tx.GetContext(ctx, &existingID, `SELECT id FROM iku_value
			WHERE label_id = ? AND periode_id = ? AND scope_type = ? AND COALESCE(scope_id,'') = COALESCE(?,'')`,
			v.LabelID, v.PeriodeID, v.ScopeType, v.ScopeID)
		if existingID == "" {
			v.ID = uuid.NewString()
			if _, err := tx.ExecContext(ctx, `INSERT INTO iku_value
				(id, label_id, periode_id, scope_type, scope_id, numerator, denominator, value, status,
				 formula_version, calculation_run_id, source_hash, data_version, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				v.ID, v.LabelID, v.PeriodeID, v.ScopeType, v.ScopeID, v.Numerator, v.Denominator,
				nullFloat(v.Value), v.Status, v.FormulaVersion, v.CalculationRunID, v.SourceHash,
				v.DataVersion, v.UpdatedAt); err != nil {
				return fmt.Errorf("insert iku_value: %w", err)
			}
		} else {
			v.ID = existingID
			if _, err := tx.ExecContext(ctx, `UPDATE iku_value SET numerator = ?, denominator = ?,
				value = ?, status = ?, formula_version = ?, calculation_run_id = ?, source_hash = ?,
				data_version = ?, updated_at = ? WHERE id = ?`,
				v.Numerator, v.Denominator, nullFloat(v.Value), v.Status, v.FormulaVersion,
				v.CalculationRunID, v.SourceHash, v.DataVersion, v.UpdatedAt, v.ID); err != nil {
				return fmt.Errorf("update iku_value: %w", err)
			}
			if _, err := tx.ExecContext(ctx, `DELETE FROM iku_value_detail WHERE iku_value_id = ?`, v.ID); err != nil {
				return err
			}
		}
		for _, d := range write.Details {
			if _, err := tx.ExecContext(ctx, `INSERT INTO iku_value_detail
				(id, iku_value_id, component_type, component_id, component_label_id, component_value,
				 numerator, denominator, weight, contribution)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				uuid.NewString(), v.ID, d.ComponentType, d.ComponentID, d.ComponentLabelID,
				nullFloat(d.ComponentValue), nullFloat(d.Numerator), nullFloat(d.Denominator),
				d.Weight, nullFloat(d.Contribution)); err != nil {
				return fmt.Errorf("insert detail: %w", err)
			}
		}
	}
	return tx.Commit()
}

func (s *Store) CreateCalculationRun(ctx context.Context, run *domain.CalculationRun) error {
	if run.ID == "" {
		run.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO calculation_run
		(id, periode_id, scope_type, scope_id, trigger, status, started_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		run.ID, run.PeriodeID, run.ScopeType, run.ScopeID, run.Trigger, run.Status, run.StartedAt)
	return err
}

func (s *Store) FinishCalculationRun(ctx context.Context, id, status, errMsg, finishedAt string) error {
	var msg any
	if errMsg != "" {
		msg = errMsg
	}
	_, err := s.db.ExecContext(ctx, `UPDATE calculation_run SET status = ?, finished_at = ?, error = ? WHERE id = ?`,
		status, finishedAt, msg, id)
	return err
}

func (s *Store) LatestCalculationRun(ctx context.Context, periodeID string) (*domain.CalculationRun, error) {
	var out domain.CalculationRun
	err := s.db.GetContext(ctx, &out, `SELECT id, periode_id, scope_type, scope_id, trigger, status,
		started_at, finished_at, error FROM calculation_run
		WHERE periode_id = ? ORDER BY started_at DESC LIMIT 1`, periodeID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

func (s *Store) ListFinalFileHashes(ctx context.Context, periodeID, scopeType, scopeID string) ([]string, error) {
	q := `SELECT DISTINCT s.file_hash FROM submission s
		WHERE s.status = 'final' AND s.periode_id = ?`
	args := []any{periodeID}
	if scopeType == "FASKES" && scopeID != "" {
		q += ` AND s.faskes_id = ?`
		args = append(args, scopeID)
	}
	q += ` ORDER BY s.file_hash`
	var out []string
	err := s.db.SelectContext(ctx, &out, q, args...)
	return out, err
}

func (s *Store) GetIndicatorValueBySubmissionRow(ctx context.Context, submissionRowID string) (*domain.IndicatorValue, error) {
	var out domain.IndicatorValue
	err := s.db.GetContext(ctx, &out, `SELECT id, submission_row_id, indikator_id, periode_id, faskes_id,
		scope_type, scope_id, numerator, denominator, value, status, data_version, created_at
		FROM indicator_value WHERE submission_row_id = ?`, submissionRowID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

func (s *Store) latestIndicatorValue(ctx context.Context, indikatorID, periodeID, scopeType, scopeID string) (*domain.IndicatorValue, error) {
	q := `SELECT iv.id, iv.submission_row_id, iv.indikator_id, iv.periode_id, iv.faskes_id,
		iv.scope_type, iv.scope_id, iv.numerator, iv.denominator, iv.value, iv.status,
		iv.data_version, iv.created_at
		FROM indicator_value iv
		JOIN submission_row sr ON sr.id = iv.submission_row_id
		JOIN submission s ON s.id = sr.submission_id
		WHERE s.status = 'final' AND iv.indikator_id = ? AND iv.periode_id = ?`
	args := []any{indikatorID, periodeID}
	if scopeType == "FASKES" && scopeID != "" {
		q += ` AND iv.faskes_id = ?`
		args = append(args, scopeID)
	}
	q += ` ORDER BY iv.created_at DESC LIMIT 1`
	var out domain.IndicatorValue
	err := s.db.GetContext(ctx, &out, q, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

type baselineJSON struct {
	CapaianTW1 *float64 `json:"capaian_tw1"`
}

func (s *Store) DashboardSPM(ctx context.Context, f domain.DashboardFilter) ([]domain.DashboardSPMItem, error) {
	rows, err := s.ListDashboardRows(ctx)
	if err != nil {
		return nil, err
	}
	baselineYear := 0
	var periodYear int
	if f.PeriodeID != "" {
		_ = s.db.GetContext(ctx, &periodYear, `SELECT tahun FROM periode WHERE id = ?`, f.PeriodeID)
	}
	baselineByLabel := map[string]float64{}
	baselineByIndicator := map[string]float64{}
	if periodYear > 0 {
		baselineYear = periodYear - 1
		baselineScopeID := f.ScopeID
		if f.ScopeType != "KABUPATEN" || baselineScopeID == "" {
			_ = s.db.GetContext(ctx, &baselineScopeID, `SELECT id FROM wilayah
				WHERE level = 'kabupaten' ORDER BY id LIMIT 1`)
		}
		baselines, err := s.ListCapaianTahunan(ctx, baselineYear, "KABUPATEN", baselineScopeID)
		if err != nil {
			return nil, err
		}
		for _, c := range baselines {
			if c.LabelID != nil {
				baselineByLabel[*c.LabelID] = c.Nilai
			}
			if c.IndikatorID != nil {
				baselineByIndicator[*c.IndikatorID] = c.Nilai
			}
		}
	}
	out := make([]domain.DashboardSPMItem, 0, len(rows))
	for _, row := range rows {
		item := domain.DashboardSPMItem{
			Kode:         row.Kode,
			LabelID:      row.LabelID,
			IndikatorID:  row.IndikatorID,
			Target:       row.Target,
			BaselineYear: baselineYear,
		}
		if row.LabelID != nil {
			if lbl, err := s.GetLabelByID(ctx, *row.LabelID); err == nil {
				item.LabelKode = &lbl.Kode
				item.LabelNama = &lbl.Nama
			}
		}
		if row.IndikatorID != nil {
			if ind, err := s.GetIndicatorByID(ctx, *row.IndikatorID); err == nil {
				item.Indikator = ind.Nama
			}
		} else {
			item.Indikator = row.Kode
		}
		var b baselineJSON
		if row.BaselineJSON != nil {
			_ = json.Unmarshal([]byte(*row.BaselineJSON), &b)
		}
		if row.LabelID != nil {
			if v, ok := baselineByLabel[*row.LabelID]; ok {
				value := v
				item.BaselinePrev = &value
			}
		}
		if item.BaselinePrev == nil && row.IndikatorID != nil {
			if v, ok := baselineByIndicator[*row.IndikatorID]; ok {
				value := v
				item.BaselinePrev = &value
			}
		}
		item.BaselineTW1 = b.CapaianTW1
		if row.LabelID != nil {
			if v, err := s.GetIKUValue(ctx, *row.LabelID, f.PeriodeID, f.ScopeType, f.ScopeID); err == nil && v != nil {
				item.IKUValueID = &v.ID
				item.LabelKode = v.LabelKode
				item.LabelNama = v.LabelNama
				item.Value = v.Value
				item.Numerator = &v.Numerator
				item.Denominator = &v.Denominator
				item.ValueStatus = string(v.Status)
				item.SourceHash = v.SourceHash
				item.DataVersion = v.DataVersion
				item.CalculationRunID = v.CalculationRunID
				item.UpdatedAt = &v.UpdatedAt
			}
		}
		if item.ValueStatus == "" {
			item.ValueStatus = string(domain.ValueNoData)
		}
		if row.IndikatorID != nil {
			if iv, err := s.latestIndicatorValue(ctx, *row.IndikatorID, f.PeriodeID, f.ScopeType, f.ScopeID); err == nil && iv != nil {
				item.IndicatorValueID = &iv.ID
			}
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *Store) TraceSPM(ctx context.Context, kode string, f domain.DashboardFilter) (*domain.SPMTrace, error) {
	row, err := s.GetDashboardRowByKode(ctx, kode)
	if err != nil {
		return nil, err
	}
	items, err := s.DashboardSPM(ctx, f)
	if err != nil {
		return nil, err
	}
	var item *domain.DashboardSPMItem
	for i := range items {
		if items[i].Kode == kode {
			item = &items[i]
			break
		}
	}
	if item == nil {
		return nil, domain.NotFound("baris dashboard")
	}
	trace := &domain.SPMTrace{Row: *item}
	if row.LabelID == nil {
		return trace, nil
	}
	iku, err := s.GetIKUValue(ctx, *row.LabelID, f.PeriodeID, f.ScopeType, f.ScopeID)
	if err != nil {
		return nil, err
	}
	if iku == nil {
		return trace, nil
	}
	trace.IKUValue = iku
	details, err := s.ListIKUValueDetails(ctx, iku.ID)
	if err != nil {
		return nil, err
	}
	trace.Details = details
	// Trace satu komponen indikator pertama yang punya nilai.
	for _, d := range details {
		if d.ComponentType != "INDIKATOR" {
			continue
		}
		iv, err := s.latestIndicatorValue(ctx, d.ComponentID, f.PeriodeID, f.ScopeType, f.ScopeID)
		if err != nil || iv == nil {
			continue
		}
		trace.IndicatorValue = iv
		sr, err := s.GetRow(ctx, iv.SubmissionRowID)
		if err != nil {
			continue
		}
		trace.SubmissionRow = sr
		sub, err := s.GetSubmission(ctx, sr.SubmissionID)
		if err == nil {
			trace.Submission = sub
			var file domain.FileObject
			if err := s.db.GetContext(ctx, &file, `SELECT id, original_name, sha256, path, mime, size,
				uploaded_by, scan_status, created_at FROM file_object WHERE id = ?`, sub.FileID); err == nil {
				trace.File = &file
			}
			reviews, _ := s.ListReviews(ctx, sub.ID)
			trace.Reviews = reviews
		}
		break
	}
	return trace, nil
}
