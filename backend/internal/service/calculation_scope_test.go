package service_test

import (
	"context"
	"math"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
)

func TestCalculateScopeHierarchy(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	if seeded, err := store.SeedDemo(ctx); err != nil || !seeded {
		t.Fatalf("seed: %v seeded=%v", err, seeded)
	}

	// A final submission for June with two village rows. This exercises the
	// lowest-grain aggregation path without changing the May seed fixture.
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO file_object
		(id, original_name, sha256, path, mime, size, uploaded_by, scan_status, created_at)
		VALUES ('file-june', 'june.xlsx', 'junehash', '/tmp/june.xlsx', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 1, 'usr-siti', 'clean', '2026-06-10T08:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission
		(id, kode, tipe, uploader_id, faskes_id, periode_id, file_id, file_hash,
		 status, step, master_release_id, workflow_version_id, created_at, updated_at)
		VALUES ('sub-june', 'S-2026-06-001', 'BNBA', 'usr-siti', 'faskes-pkm-jereweh',
		 'per-2026-06', 'file-june', 'junehash', 'final', 'FINAL', 'mr-2026', 'wf-default-v1',
		 '2026-06-10T08:00:00Z', '2026-06-10T10:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission_revision
		(id, submission_id, revision_no, parent_revision_id, file_id, created_by, created_at)
		VALUES ('rev-june', 'sub-june', 1, NULL, 'file-june', 'usr-siti', '2026-06-10T08:00:00Z')`); err != nil {
		t.Fatal(err)
	}

	rows := []struct {
		rowID, desaID, desaKode string
		num, den                float64
	}{
		{"row-june-1", "wil-desa-dalam", "5201.01.01", 92, 120},
		{"row-june-2", "wil-desa-beru", "5201.01.02", 70, 100},
	}
	for i, r := range rows {
		rowNo := i + 1
		if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission_row
			(id, submission_id, group_id, row_no, kode_iku_raw, kode_indikator_raw,
			 matched_indikator_id, matched_label_id, scope_type, scope_code, nilai,
			 numerator, denominator, classification_status, error_message, created_at)
			VALUES (?, 'sub-june', NULL, ?, 'IKU-2026-001.01', 'IND-001',
			 'ind-001', 'lbl-iku-001-01', 'DESA', ?, ?, ?, ?, 'MATCHED', NULL,
			 '2026-06-10T08:01:00Z')`,
			r.rowID, rowNo, r.desaKode, r.num/r.den*100, r.num, r.den); err != nil {
			t.Fatal(err)
		}
		if _, err := store.DB().ExecContext(ctx, `INSERT INTO indicator_value
			(id, submission_row_id, indikator_id, periode_id, faskes_id, scope_type,
			 scope_id, numerator, denominator, value, status, data_version, created_at)
			VALUES (?, ?, 'ind-001', 'per-2026-06', 'faskes-pkm-jereweh', 'DESA',
			 ?, ?, ?, ?, 'COMPLETE', 'vtest', '2026-06-10T08:01:00Z')`,
			"iv-"+r.rowID, r.rowID, r.desaID, r.num, r.den, r.num/r.den*100); err != nil {
			t.Fatal(err)
		}
	}

	calc := service.NewCalculationService(store, store)
	if _, err := calc.CalculatePeriod(ctx, "per-2026-06", "test"); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		scopeType, scopeID string
		want               float64
	}{
		{"DESA", "wil-desa-dalam", 92.0 / 120.0 * 100},
		{"DESA", "wil-desa-beru", 70.0 / 100.0 * 100},
		{"KECAMATAN", "wil-kec-taliwang", (92.0 + 70.0) / (120.0 + 100.0) * 100},
		{"KABUPATEN", "wil-kab-swb", (92.0 + 70.0) / (120.0 + 100.0) * 100},
		{"FASKES", "faskes-pkm-taliwang", (92.0 + 70.0) / (120.0 + 100.0) * 100},
	}
	for _, c := range cases {
		v, err := store.GetIKUValue(ctx, "lbl-iku-001-01", "per-2026-06", c.scopeType, c.scopeID)
		if err != nil || v == nil || v.Value == nil {
			t.Fatalf("%s/%s: value=%v err=%v", c.scopeType, c.scopeID, v, err)
		}
		if math.Abs(*v.Value-c.want) > 0.0001 {
			t.Fatalf("%s/%s: got %.6f want %.6f", c.scopeType, c.scopeID, *v.Value, c.want)
		}
	}

	jereweh, err := store.GetIKUValue(ctx, "lbl-iku-001-01", "per-2026-06", "FASKES", "faskes-pkm-jereweh")
	if err != nil {
		t.Fatal(err)
	}
	if jereweh != nil && jereweh.Value != nil {
		t.Fatalf("desa Taliwang tidak boleh masuk faskes Jereweh: %+v", jereweh)
	}

	parent, err := store.GetIKUValue(ctx, "lbl-iku-001", "per-2026-06", "KABUPATEN", "wil-kab-swb")
	if err != nil || parent == nil || parent.Value == nil {
		t.Fatalf("parent: %+v err=%v", parent, err)
	}
	if parent.Status != "PARTIAL" {
		t.Fatalf("parent status: %s", parent.Status)
	}
	if math.Abs(*parent.Value-(92.0+70.0)/(120.0+100.0)*100) > 0.0001 {
		t.Fatalf("parent value: %.6f", *parent.Value)
	}

	hashes, err := store.ListFinalFileHashes(ctx, "per-2026-06", "FASKES", "faskes-pkm-taliwang")
	if err != nil || len(hashes) == 0 {
		t.Fatalf("hash faskes domisili: %v %v", hashes, err)
	}
	domisiliValues, err := store.ListFinalIndicatorValuesForFaskesDomisili(ctx, "per-2026-06", "faskes-pkm-taliwang")
	if err != nil || len(domisiliValues) == 0 {
		t.Fatalf("nilai faskes domisili: %v %v", domisiliValues, err)
	}
	jerewehValues, err := store.ListFinalIndicatorValuesForFaskesDomisili(ctx, "per-2026-06", "faskes-pkm-jereweh")
	if err != nil {
		t.Fatal(err)
	}
	if len(jerewehValues) != 0 {
		t.Fatalf("nilai desa Taliwang bocor ke Jereweh: %+v", jerewehValues)
	}
	hashes, err = store.ListFinalFileHashes(ctx, "per-2026-06", "KECAMATAN", "wil-kec-taliwang")
	if err != nil || len(hashes) == 0 {
		t.Fatalf("hash kecamatan: %v %v", hashes, err)
	}

	dash := service.NewDashboardService(store, store, store)
	items, _, err := dash.SPM(ctx, domain.DashboardFilter{
		PeriodeID: "per-2026-06", ScopeType: "KECAMATAN", ScopeID: "wil-kec-taliwang",
	})
	if err != nil {
		t.Fatal(err)
	}
	var got float64
	for _, item := range items {
		if item.LabelID != nil && *item.LabelID == "lbl-iku-001-01" && item.Value != nil {
			got = *item.Value
		}
	}
	if math.Abs(got-(92.0+70.0)/(120.0+100.0)*100) > 0.0001 {
		t.Fatalf("dashboard kecamatan: %.6f", got)
	}

	detail, err := dash.Detail(ctx, "01", domain.DashboardFilter{
		PeriodeID: "per-2026-06", ScopeType: "KECAMATAN", ScopeID: "wil-kec-taliwang",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(detail.Wilayah) == 0 {
		t.Fatal("detail wilayah kosong")
	}
}
