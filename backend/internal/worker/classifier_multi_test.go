package worker_test

import (
	"context"
	"path/filepath"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository/sqlite"
	"matasidik/backend/internal/worker"
)

func TestClassifyMultipleIKUCodes(t *testing.T) {
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

	// Add a second label and link it to the same indicator. This simulates the
	// many-to-many labels from sheet 2 (IKU + IKK) without changing the seed.
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO label
		(id, kode, nama, tipe, category_id, parent_id, master_release_id,
		 aggregation_method, allow_direct_input, is_rollup_only, level, path, aktif, berlaku_dari)
		VALUES ('lbl-ikk-test', 'IKK-TEST', 'IKK Test', 'IKK', 'cat-ikk', NULL, 'mr-2026',
		 'RATIO_OF_SUMS', 1, 0, 2, 'IKK-TEST', 1, '2026-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO indikator_label
		(indikator_id, label_id, bobot, is_primary) VALUES ('ind-001', 'lbl-ikk-test', 1, 0)`); err != nil {
		t.Fatal(err)
	}

	sub, err := store.GetSubmission(ctx, "sub-demo-001")
	if err != nil {
		t.Fatal(err)
	}
	num, den := 92.0, 120.0
	raw := []worker.RawRow{{
		RowNo: 1, Periode: "2026-05",
		KodeIKU:       "IKU-2026-001.01; IKK-TEST",
		KodeIndikator: "IND-001", ScopeType: "DESA", ScopeCode: "5201.01.01",
		Numerator: &num, Denominator: &den,
	}}
	classified, err := worker.NewClassifier(store).Classify(ctx, sub, raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Rows) != 1 {
		t.Fatalf("hasil klasifikasi: %d", len(classified.Rows))
	}
	row := classified.Rows[0]
	if row.Classification != domain.ClassMatched {
		t.Fatalf("klasifikasi: got %s (%s)", row.Classification, row.ErrorMessage)
	}
	if row.MatchedLabelID != "lbl-iku-001-01" {
		t.Fatalf("label utama: got %v", row.MatchedLabelID)
	}
	if len(classified.Values) != 1 {
		t.Fatalf("nilai indikator: %d", len(classified.Values))
	}
}
