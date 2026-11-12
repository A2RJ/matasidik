package worker_test

import (
	"context"
	"path/filepath"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository/sqlite"
	"matasidik/backend/internal/worker"
)

func TestParseAndClassifySample(t *testing.T) {
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
	sub, err := store.GetSubmission(ctx, "sub-demo-001")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := worker.ParseFile("../repository/sqlite/testdata/bnba-2026-05.csv")
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) != 3 {
		t.Fatalf("jumlah baris: %d", len(raw))
	}
	classified, err := worker.NewClassifier(store).Classify(ctx, sub, raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(classified.Rows) != 3 {
		t.Fatalf("hasil klasifikasi: %d", len(classified.Rows))
	}
	got := []domain.ClassificationStatus{
		classified.Rows[0].Classification,
		classified.Rows[1].Classification,
		classified.Rows[2].Classification,
	}
	want := []domain.ClassificationStatus{domain.ClassMatched, domain.ClassMatched, domain.ClassRollupOnly}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("baris %d: got %s want %s", i+1, got[i], want[i])
		}
	}
	if len(classified.Values) != 2 {
		t.Fatalf("nilai indikator matched: %d", len(classified.Values))
	}
}
