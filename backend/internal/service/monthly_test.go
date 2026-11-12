package service_test

import (
	"context"
	"math"
	"testing"

	"matasidik/backend/internal/service"
)

func TestMonthlySPMReport(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	if _, err := store.SeedDemo(ctx); err != nil {
		t.Fatal(err)
	}
	calc := service.NewCalculationService(store, store)
	if _, err := calc.CalculatePeriod(ctx, "per-2026-05", "test"); err != nil {
		t.Fatal(err)
	}
	dashboard := service.NewDashboardService(store, store, store)
	report, err := dashboard.MonthlySPM(ctx, 2026, domainFilter("per-2026-05"))
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Rows) != 12 {
		t.Fatalf("jumlah baris laporan: %d", len(report.Rows))
	}
	var ibuHamil *struct {
		targetTahunan float64
		targetBulanan float64
		may           float64
		mayStatus     string
		janStatus     string
	}
	for _, row := range report.Rows {
		if row.Kode != "01" {
			continue
		}
		if row.TargetTahunan == nil || row.TargetBulanan == nil {
			t.Fatalf("target layanan ibu hamil belum ada: %+v", row)
		}
		if math.Abs(*row.TargetTahunan-2000) > 0.001 || math.Abs(*row.TargetBulanan-166.6666667) > 0.001 {
			t.Fatalf("target layanan ibu hamil salah: %+v", row)
		}
		if len(row.Months) != 12 {
			t.Fatalf("jumlah bulan: %d", len(row.Months))
		}
		if row.Months[0].Status != "NO_DATA" {
			t.Fatalf("Januari harus NO_DATA: %+v", row.Months[0])
		}
		if row.Months[4].Numerator == nil || math.Abs(*row.Months[4].Numerator-450) > 0.001 {
			t.Fatalf("numerator Mei salah: %+v", row.Months[4])
		}
		if row.Months[4].Status != "ON_TARGET" {
			t.Fatalf("status Mei salah: %s", row.Months[4].Status)
		}
		ibuHamil = &struct {
			targetTahunan float64
			targetBulanan float64
			may           float64
			mayStatus     string
			janStatus     string
		}{*row.TargetTahunan, *row.TargetBulanan, *row.Months[4].Numerator, row.Months[4].Status, row.Months[0].Status}
	}
	if ibuHamil == nil {
		t.Fatal("baris ibu hamil tidak ditemukan")
	}
}
