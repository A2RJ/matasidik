package service_test

import (
	"context"
	"math"
	"path/filepath"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository/sqlite"
	"matasidik/backend/internal/service"
)

func domainFilter(period string) domain.DashboardFilter {
	return domain.DashboardFilter{PeriodeID: period}
}

func newTestStore(t *testing.T) *sqlite.Store {
	t.Helper()
	dir := t.TempDir()
	dsn := "file:" + filepath.Join(dir, "test.db") + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	store, err := sqlite.New(dsn, filepath.Join(dir, "files"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func TestCalculateAndDashboard(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	seeded, err := store.SeedDemo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !seeded {
		t.Fatal("seed demo tidak berjalan")
	}
	calc := service.NewCalculationService(store, store)
	if _, err := calc.CalculatePeriod(ctx, "per-2026-05", "test"); err != nil {
		t.Fatal(err)
	}
	values, err := store.ListIKUValues(ctx, "per-2026-05", "KABUPATEN", "wil-kab-swb")
	if err != nil {
		t.Fatal(err)
	}
	found := map[string]float64{}
	for _, v := range values {
		if v.Value != nil {
			found[*v.LabelKode] = *v.Value
		}
	}
	if math.Abs(found["IKU-2026-001"]-92.3913043478261) > 0.0001 {
		t.Fatalf("nilai induk salah: %v", found["IKU-2026-001"])
	}
	if found["IKU-2026-001.01"] != 90 || found["IKU-2026-001.02"] != 95 || found["IKU-2026-001.03"] != 100 {
		t.Fatalf("nilai anak salah: %+v", found)
	}
	dash := service.NewDashboardService(store, store, store)
	items, summary, err := dash.SPM(ctx, domainFilter("per-2026-05"))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 12 {
		t.Fatalf("jumlah baris dashboard: %d", len(items))
	}
	if summary.ParentValue == nil || math.Abs(*summary.ParentValue-92.3913043478261) > 0.0001 {
		t.Fatalf("summary induk: %+v", summary)
	}
	trace, err := dash.Trace(ctx, "01", domainFilter("per-2026-05"))
	if err != nil {
		t.Fatal(err)
	}
	if trace.File == nil || trace.SubmissionRow == nil || trace.IKUValue == nil {
		t.Fatalf("trace tidak lengkap: %+v", trace)
	}
	if trace.File.SHA256 == "" || trace.SubmissionRow.RowNo != 1 {
		t.Fatalf("trace file/baris salah: %+v", trace)
	}
	detail, err := dash.Detail(ctx, "01", domainFilter("per-2026-05"))
	if err != nil {
		t.Fatal(err)
	}
	if len(detail.Faskes) != 2 || detail.Faskes[0].Uploaded == detail.Faskes[1].Uploaded {
		t.Fatalf("coverage faskes salah: %+v", detail.Faskes)
	}
	if len(detail.Provenance) == 0 {
		t.Fatal("provenance kosong")
	}
	chart, err := dash.IKUByFaskes(ctx, domainFilter("per-2026-05"))
	if err != nil {
		t.Fatal(err)
	}
	if len(chart.Faskes) != 2 || len(chart.Series) < 3 {
		t.Fatalf("chart tidak lengkap: %+v", chart)
	}
	export, err := dash.ExportSPM(ctx, domainFilter("per-2026-05"), "xlsx")
	if err != nil || len(export.Data) == 0 {
		t.Fatalf("export xlsx gagal: %v", err)
	}
}
