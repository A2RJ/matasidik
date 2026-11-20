package service_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
	"matasidik/backend/internal/worker"
)

func TestUploadGeneralTemplateRatioAndCount(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	if seeded, err := store.SeedDemo(ctx); err != nil || !seeded {
		t.Fatalf("seed: %v seeded=%v", err, seeded)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO label
		(id, kode, nama, tipe, category_id, parent_id, master_release_id, aggregation_method,
		 allow_direct_input, is_rollup_only, pj_bidang_id, level, path, aktif, berlaku_dari)
		VALUES ('lbl-count-monthly', 'IKU-COUNT-MONTHLY', 'Uji Count Bulanan', 'IKU', 'cat-iku',
		 NULL, 'mr-2026', 'SUM', 1, 0, 'bid-kesmas', 1, 'IKU-COUNT-MONTHLY', 1, '2026-01-01')`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO indikator
		(id, kode, nama, satuan, tipe_nilai, sifat_input, frekuensi, pj_bidang_id, aktif)
		VALUES ('ind-count-monthly', 'IND-COUNT-MONTHLY', 'Jumlah kasus bulanan', 'Kasus',
		 'COUNT', 'puskesmas', 'bulanan', 'bid-kesmas', 1)`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO indikator_label
		(indikator_id, label_id, bobot, is_primary)
		VALUES ('ind-count-monthly', 'lbl-count-monthly', 1, 1)`); err != nil {
		t.Fatal(err)
	}

	f := excelize.NewFile()
	_ = f.SetSheetName("Sheet1", "Petunjuk")
	_, _ = f.NewSheet("Rasio")
	_, _ = f.NewSheet("Nilai")
	rasioHeaders := []any{"periode", "kode_faskes", "kode_iku", "kode_indikator", "scope_type", "scope_code", "numerator", "denominator"}
	_ = f.SetSheetRow("Rasio", "A1", &rasioHeaders)
	rasioRow := []any{"2026-05", "PKM-TALIWANG", "IKU-2026-001.01", "IND-001", "DESA", "5201.01.01", 92, 120}
	_ = f.SetSheetRow("Rasio", "A2", &rasioRow)
	nilaiHeaders := []any{"periode", "kode_faskes", "kode_iku", "kode_indikator", "scope_type", "scope_code", "nilai"}
	_ = f.SetSheetRow("Nilai", "A1", &nilaiHeaders)
	nilaiRow := []any{"2026-05", "PKM-TALIWANG", "IKU-COUNT-MONTHLY", "IND-COUNT-MONTHLY", "KABUPATEN", "5201", 7}
	_ = f.SetSheetRow("Nilai", "A2", &nilaiRow)
	path := filepath.Join(t.TempDir(), "upload.xlsx")
	if err := f.SaveAs(path); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	audit := service.NewAuditService(store)
	calc := service.NewCalculationService(store, store)
	subs := service.NewSubmissionService(store, store, audit, calc, t.TempDir(), 15<<20)
	faskesID := "faskes-pkm-taliwang"
	sub, job, err := subs.Upload(ctx, &domain.User{
		ID: "usr-siti", Role: domain.RolePKMOperator, FaskesID: &faskesID,
	}, service.UploadInput{
		File: file, Filename: "upload.xlsx", Size: info.Size(), PeriodeKode: "2026-05",
		FaskesID: faskesID, IdempotencyKey: "test-general-template",
	})
	if err != nil {
		t.Fatal(err)
	}
	if sub == nil || job == nil {
		t.Fatalf("submission/job kosong: %+v %+v", sub, job)
	}
	w := worker.New("test-worker", store, store, store)
	if _, err := w.ProcessOnce(ctx); err != nil {
		t.Fatal(err)
	}
	rows, err := store.GetRows(ctx, sub.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows=%d", len(rows))
	}
	got := map[string]float64{}
	for _, row := range rows {
		if row.Classification != domain.ClassMatched {
			t.Fatalf("row %d: %s (%v)", row.RowNo, row.Classification, row.ErrorMessage)
		}
		iv, err := store.GetIndicatorValueBySubmissionRow(ctx, row.ID)
		if err != nil || iv == nil {
			t.Fatalf("indicator value row %d: %+v err=%v", row.RowNo, iv, err)
		}
		got[*row.KodeIndikatorRaw] = iv.Value
	}
	if got["IND-001"] != 92.0/120.0*100 || got["IND-COUNT-MONTHLY"] != 7 {
		t.Fatalf("nilai tersimpan salah: %+v", got)
	}
}
