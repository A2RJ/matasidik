package worker_test

import (
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"

	"matasidik/backend/internal/worker"
)

func TestParseXLSXReadsRasioAndNilaiSheets(t *testing.T) {
	f := excelize.NewFile()
	if err := f.SetSheetName("Sheet1", "Petunjuk"); err != nil {
		t.Fatal(err)
	}
	if _, err := f.NewSheet("Rasio"); err != nil {
		t.Fatal(err)
	}
	if _, err := f.NewSheet("Nilai"); err != nil {
		t.Fatal(err)
	}
	if _, err := f.NewSheet("Referensi"); err != nil {
		t.Fatal(err)
	}
	headers := []any{"periode", "kode_faskes", "kode_iku", "kode_indikator", "scope_type", "scope_code", "numerator", "denominator"}
	if err := f.SetSheetRow("Rasio", "A1", &headers); err != nil {
		t.Fatal(err)
	}
	ratio := []any{"2026-05", "PKM-TALIWANG", "IKU-2026-001.01", "IND-001", "DESA", "5201.01.01", 92, 120}
	if err := f.SetSheetRow("Rasio", "A2", &ratio); err != nil {
		t.Fatal(err)
	}
	nilaiHeaders := []any{"periode", "kode_faskes", "kode_iku", "kode_indikator", "scope_type", "scope_code", "nilai"}
	if err := f.SetSheetRow("Nilai", "A1", &nilaiHeaders); err != nil {
		t.Fatal(err)
	}
	scalar := []any{"2026", "PKM-TALIWANG", "IKU-TEST", "IND-TEST", "KABUPATEN", "5201", 73.5}
	if err := f.SetSheetRow("Nilai", "A2", &scalar); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "template.xlsx")
	if err := f.SaveAs(path); err != nil {
		t.Fatal(err)
	}

	rows, err := worker.ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows=%d", len(rows))
	}
	if rows[0].RowNo != 1 || rows[0].Numerator == nil || *rows[0].Numerator != 92 || rows[0].Nilai != nil {
		t.Fatalf("baris rasio salah: %+v", rows[0])
	}
	if rows[1].RowNo != 2 || rows[1].Nilai == nil || *rows[1].Nilai != 73.5 || rows[1].Numerator != nil {
		t.Fatalf("baris nilai salah: %+v", rows[1])
	}
}
