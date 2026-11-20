package worker_test

import (
	"context"
	"path/filepath"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository/sqlite"
	"matasidik/backend/internal/worker"
)

func TestClassifyCountAndScalarWithoutDenominator(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	store, err := sqlite.New("file:"+filepath.Join(dir, "test.db")+"?_pragma=foreign_keys(1)", filepath.Join(dir, "files"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if _, err := store.SeedDemo(ctx); err != nil {
		t.Fatal(err)
	}
	// Seed hanya berisi indikator rasio. Ubah dua indikator menjadi count dan
	// scalar untuk menguji jalur tanpa denominator.
	if _, err := store.DB().ExecContext(ctx, `UPDATE indikator SET tipe_nilai = 'COUNT', satuan = 'Kasus' WHERE id = 'ind-001'`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `UPDATE indikator SET tipe_nilai = 'SCALAR', satuan = '(tahun)' WHERE id = 'ind-002'`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `UPDATE indikator SET tipe_nilai = 'RATIO', satuan = 'Rasio' WHERE id = 'ind-003'`); err != nil {
		t.Fatal(err)
	}
	sub, err := store.GetSubmission(ctx, "sub-demo-001")
	if err != nil {
		t.Fatal(err)
	}
	countValue, scalarValue := 7.0, 73.5
	ratioNum, ratioDen := 4.0, 2.0
	raw := []worker.RawRow{
		{
			RowNo: 1, Periode: "2026-05", KodeIKU: "IKU-2026-001.01",
			KodeIndikator: "IND-001", ScopeType: "DESA", ScopeCode: "5201.01.01",
			Nilai: &countValue,
		},
		{
			RowNo: 2, Periode: "2026-05", KodeIKU: "IKU-2026-001.02",
			KodeIndikator: "IND-002", ScopeType: "KABUPATEN", ScopeCode: "5201",
			Nilai: &scalarValue,
		},
		{
			RowNo: 3, Periode: "2026-05", KodeIKU: "IKU-2026-001.03",
			KodeIndikator: "IND-003", ScopeType: "KABUPATEN", ScopeCode: "5201",
			Numerator: &ratioNum, Denominator: &ratioDen,
		},
	}
	classified, err := worker.NewClassifier(store).Classify(ctx, sub, raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Rows) != 3 || len(classified.Values) != 3 {
		t.Fatalf("rows=%d values=%d", len(classified.Rows), len(classified.Values))
	}
	for i, row := range classified.Rows {
		if row.Classification != domain.ClassMatched {
			t.Fatalf("baris %d: %s (%v)", i+1, row.Classification, row.ErrorMessage)
		}
		if i < 2 && (row.Denominator != 0 || row.Numerator != 0) {
			t.Fatalf("baris %d masih memakai numerator/denominator: %+v", i+1, row)
		}
	}
	if classified.Rows[2].Numerator != 4 || classified.Rows[2].Denominator != 2 {
		t.Fatalf("rasio non-persen kehilangan num/den: %+v", classified.Rows[2])
	}
	if classified.Rows[0].Value != 7 || classified.Rows[1].Value != 73.5 || classified.Rows[2].Value != 2 {
		t.Fatalf("nilai count/scalar/rasio salah: %+v", classified.Rows)
	}
}

func TestClassifyAnnualIndicatorRequiresAnnualPeriod(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	store, err := sqlite.New("file:"+filepath.Join(dir, "test.db")+"?_pragma=foreign_keys(1)", filepath.Join(dir, "files"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if _, err := store.SeedDemo(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `UPDATE indikator SET tipe_nilai = 'SCALAR', satuan = '(tahun)', frekuensi = 'tahunan' WHERE id = 'ind-001'`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO file_object
		(id, original_name, sha256, path, mime, size, uploaded_by, scan_status, created_at)
		VALUES ('file-annual-test', 'annual.xlsx', 'annualtest', '/tmp/annual-test.xlsx', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 1, 'usr-siti', 'clean', '2026-12-01T08:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission
		(id, kode, tipe, uploader_id, faskes_id, periode_id, file_id, file_hash,
		 status, step, master_release_id, workflow_version_id, created_at, updated_at)
		VALUES ('sub-annual-test', 'S-2026-ANNUAL-TEST', 'BNBA', 'usr-siti', 'faskes-pkm-taliwang',
		 'per-2026-annual', 'file-annual-test', 'annualtest', 'draf', 'DRAFT_PKM', 'mr-2026', 'wf-default-v2',
		 '2026-12-01T08:00:00Z', '2026-12-01T08:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	annualSub, err := store.GetSubmission(ctx, "sub-annual-test")
	if err != nil {
		t.Fatal(err)
	}
	value := 73.5
	raw := []worker.RawRow{{
		RowNo: 1, Periode: "2026", KodeIKU: "IKU-2026-001.01",
		KodeIndikator: "IND-001", ScopeType: "KABUPATEN", ScopeCode: "5201", Nilai: &value,
	}}
	classified, err := worker.NewClassifier(store).Classify(ctx, annualSub, raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Rows) != 1 || classified.Rows[0].Classification != domain.ClassMatched {
		t.Fatalf("annual scalar tidak matched: %+v", classified.Rows)
	}

	monthlySub, err := store.GetSubmission(ctx, "sub-demo-001")
	if err != nil {
		t.Fatal(err)
	}
	raw[0].Periode = "2026-05"
	classified, err = worker.NewClassifier(store).Classify(ctx, monthlySub, raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Rows) != 1 || classified.Rows[0].Classification != domain.ClassInvalidValue {
		t.Fatalf("annual indicator di periode bulanan seharusnya ditolak: %+v", classified.Rows)
	}
}
