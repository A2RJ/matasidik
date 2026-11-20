package service_test

import (
	"context"
	"math"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
)

// TestCalculateCountScalarAnnual memastikan count/scalar tidak membutuhkan
// denominator dan periode tahunan (bulan=0) ikut dihitung.
func TestCalculateCountScalarAnnual(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	if seeded, err := store.SeedDemo(ctx); err != nil || !seeded {
		t.Fatalf("seed: %v seeded=%v", err, seeded)
	}

	// Label count dan scalar terpisah agar agregasi SUM dan LAST_VALUE teruji.
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO label
		(id, kode, nama, tipe, category_id, parent_id, master_release_id, aggregation_method,
		 allow_direct_input, is_rollup_only, pj_bidang_id, level, path, aktif, berlaku_dari)
		VALUES
		('lbl-count-test', 'IKU-TEST-COUNT', 'Uji Count', 'IKU', 'cat-iku', NULL, 'mr-2026', 'SUM',
		 1, 0, 'bid-kesmas', 1, 'IKU-TEST-COUNT', 1, '2026-01-01'),
		('lbl-scalar-test', 'IKU-TEST-SCALAR', 'Uji Scalar', 'IKU', 'cat-iku', NULL, 'mr-2026', 'SUM',
		 1, 0, 'bid-kesmas', 1, 'IKU-TEST-SCALAR', 1, '2026-01-01'),
		('lbl-rasio-test', 'IKU-TEST-RASIO', 'Uji Rasio Non-Persen', 'IKU', 'cat-iku', NULL, 'mr-2026', 'RATIO_OF_SUMS',
		 1, 0, 'bid-kesmas', 1, 'IKU-TEST-RASIO', 1, '2026-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO indikator
		(id, kode, nama, satuan, tipe_nilai, sifat_input, frekuensi, pj_bidang_id, aktif)
		VALUES
		('ind-count-test', 'IND-TEST-COUNT', 'Jumlah kasus uji', 'Kasus', 'COUNT', 'puskesmas', 'tahunan', 'bid-kesmas', 1),
		('ind-scalar-test', 'IND-TEST-SCALAR', 'Nilai scalar uji', '(tahun)', 'SCALAR', 'puskesmas', 'tahunan', 'bid-kesmas', 1),
		('ind-rasio-test', 'IND-TEST-RASIO', 'Rasio tenaga uji', 'Rasio', 'RATIO', 'puskesmas', 'tahunan', 'bid-kesmas', 1)`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO indikator_label
		(indikator_id, label_id, bobot, is_primary) VALUES
		('ind-count-test', 'lbl-count-test', 1, 1),
		('ind-scalar-test', 'lbl-scalar-test', 1, 1),
		('ind-rasio-test', 'lbl-rasio-test', 1, 1)`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO dashboard_row
		(id, kode, label_id, indikator_id, urutan, target) VALUES
		('dash-count-test', 'T-COUNT', 'lbl-count-test', 'ind-count-test', 90, 10),
		('dash-scalar-test', 'T-SCALAR', 'lbl-scalar-test', 'ind-scalar-test', 91, 75),
		('dash-rasio-test', 'T-RASIO', 'lbl-rasio-test', 'ind-rasio-test', 92, 2)`); err != nil {
		t.Fatal(err)
	}

	if _, err := store.DB().ExecContext(ctx, `INSERT INTO file_object
		(id, original_name, sha256, path, mime, size, uploaded_by, scan_status, created_at)
		VALUES ('file-annual', 'annual.xlsx', 'annualhash', '/tmp/annual.xlsx',
		 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 1, 'usr-siti', 'clean', '2026-12-01T08:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission
		(id, kode, tipe, uploader_id, faskes_id, periode_id, file_id, file_hash,
		 status, step, master_release_id, workflow_version_id, created_at, updated_at)
		VALUES ('sub-annual', 'S-2026-ANNUAL', 'BNBA', 'usr-siti', 'faskes-pkm-taliwang',
		 'per-2026-annual', 'file-annual', 'annualhash', 'final', 'FINAL', 'mr-2026', 'wf-default-v1',
		 '2026-12-01T08:00:00Z', '2026-12-01T09:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission_revision
		(id, submission_id, revision_no, parent_revision_id, file_id, created_by, created_at)
		VALUES ('rev-annual', 'sub-annual', 1, NULL, 'file-annual', 'usr-siti', '2026-12-01T08:00:00Z')`); err != nil {
		t.Fatal(err)
	}

	type rowSeed struct {
		rowID, indicatorID, labelID, code string
		value                             float64
		createdAt                         string
	}
	rows := []rowSeed{
		{"row-count", "ind-count-test", "lbl-count-test", "IKU-TEST-COUNT", 7, "2026-12-01T08:10:00Z"},
		{"row-scalar-old", "ind-scalar-test", "lbl-scalar-test", "IKU-TEST-SCALAR", 73.5, "2026-12-01T08:10:00Z"},
		{"row-scalar-new", "ind-scalar-test", "lbl-scalar-test", "IKU-TEST-SCALAR", 74.2, "2026-12-01T08:20:00Z"},
	}
	for i, r := range rows {
		if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission_row
			(id, submission_id, group_id, row_no, kode_iku_raw, kode_indikator_raw,
			 matched_indikator_id, matched_label_id, scope_type, scope_code, nilai,
			 numerator, denominator, classification_status, created_at)
			VALUES (?, 'sub-annual', NULL, ?, ?, ?, ?, ?, 'KABUPATEN', '5201', ?, 0, 0, 'MATCHED', ?)`,
			r.rowID, i+1, r.code, r.indicatorID, r.indicatorID, r.labelID, r.value, r.createdAt); err != nil {
			t.Fatal(err)
		}
		if _, err := store.DB().ExecContext(ctx, `INSERT INTO indicator_value
			(id, submission_row_id, indikator_id, periode_id, faskes_id, scope_type, scope_id,
			 numerator, denominator, value, status, data_version, created_at)
			VALUES (?, ?, ?, 'per-2026-annual', 'faskes-pkm-taliwang', 'KABUPATEN', 'wil-kab-swb',
			 ?, 0, ?, 'COMPLETE', 'vtest', ?)`,
			"iv-"+r.rowID, r.rowID, r.indicatorID, r.value, r.value, r.createdAt); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission_row
		(id, submission_id, group_id, row_no, kode_iku_raw, kode_indikator_raw,
		 matched_indikator_id, matched_label_id, scope_type, scope_code, nilai,
		 numerator, denominator, classification_status, created_at)
		VALUES ('row-rasio', 'sub-annual', NULL, 4, 'IKU-TEST-RASIO', 'IND-TEST-RASIO',
		 'ind-rasio-test', 'lbl-rasio-test', 'KABUPATEN', '5201', 200, 4, 2, 'MATCHED', '2026-12-01T08:30:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO indicator_value
		(id, submission_row_id, indikator_id, periode_id, faskes_id, scope_type, scope_id,
		 numerator, denominator, value, status, data_version, created_at)
		VALUES ('iv-row-rasio', 'row-rasio', 'ind-rasio-test', 'per-2026-annual', 'faskes-pkm-taliwang',
		 'KABUPATEN', 'wil-kab-swb', 4, 2, 200, 'COMPLETE', 'vtest', '2026-12-01T08:30:00Z')`); err != nil {
		t.Fatal(err)
	}

	calc := service.NewCalculationService(store, store)
	if _, err := calc.CalculatePeriod(ctx, "per-2026-annual", "test"); err != nil {
		t.Fatal(err)
	}
	count, err := store.GetIKUValue(ctx, "lbl-count-test", "per-2026-annual", "KABUPATEN", "wil-kab-swb")
	if err != nil || count == nil || count.Value == nil {
		t.Fatalf("count: %+v err=%v", count, err)
	}
	if *count.Value != 7 || count.Numerator != 7 || count.Denominator != 0 || count.Status != domain.ValueComplete {
		t.Fatalf("nilai count salah: %+v", count)
	}
	scalar, err := store.GetIKUValue(ctx, "lbl-scalar-test", "per-2026-annual", "KABUPATEN", "wil-kab-swb")
	if err != nil || scalar == nil || scalar.Value == nil {
		t.Fatalf("scalar: %+v err=%v", scalar, err)
	}
	if math.Abs(*scalar.Value-74.2) > 0.0001 || scalar.Denominator != 0 {
		t.Fatalf("nilai scalar bukan LAST_VALUE: %+v", scalar)
	}
	rasio, err := store.GetIKUValue(ctx, "lbl-rasio-test", "per-2026-annual", "KABUPATEN", "wil-kab-swb")
	if err != nil || rasio == nil || rasio.Value == nil {
		t.Fatalf("rasio: %+v err=%v", rasio, err)
	}
	if math.Abs(*rasio.Value-2) > 0.0001 || rasio.Numerator != 4 || rasio.Denominator != 2 {
		t.Fatalf("nilai rasio non-persen salah: %+v", rasio)
	}

	dash := service.NewDashboardService(store, store, store)
	items, _, err := dash.SPM(ctx, domain.DashboardFilter{
		PeriodeID: "per-2026-annual", ScopeType: "KABUPATEN", ScopeID: "wil-kab-swb",
	})
	if err != nil {
		t.Fatal(err)
	}
	found := map[string]domain.DashboardSPMItem{}
	for _, item := range items {
		found[item.Kode] = item
	}
	countItem := found["T-COUNT"]
	if countItem.TipeNilai != string(domain.TipeCount) || countItem.Satuan != "Kasus" || countItem.Value == nil || *countItem.Value != 7 || countItem.DisplayStatus != "BELOW_TARGET" {
		t.Fatalf("dashboard count salah: %+v", countItem)
	}
	scalarItem := found["T-SCALAR"]
	if scalarItem.TipeNilai != string(domain.TipeScalar) || scalarItem.Value == nil || math.Abs(*scalarItem.Value-74.2) > 0.0001 {
		t.Fatalf("dashboard scalar salah: %+v", scalarItem)
	}
	rasioItem := found["T-RASIO"]
	if rasioItem.TipeNilai != string(domain.TipeRatio) || rasioItem.Satuan != "Rasio" || rasioItem.Value == nil || math.Abs(*rasioItem.Value-2) > 0.0001 {
		t.Fatalf("dashboard rasio non-persen salah: %+v", rasioItem)
	}
	if _, _, err := dash.SPM(ctx, domain.DashboardFilter{PeriodeID: "2026", ScopeType: "KABUPATEN", ScopeID: "wil-kab-swb"}); err != nil {
		t.Fatalf("filter periode tahunan gagal: %v", err)
	}
}
