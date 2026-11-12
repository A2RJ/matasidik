package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/xuri/excelize/v2"

	"matasidik/backend/internal/domain"
)

type ExportFile struct {
	Data        []byte
	ContentType string
	Filename    string
}

func (s *DashboardService) ExportSPM(ctx context.Context, f domain.DashboardFilter, format string) (*ExportFile, error) {
	items, summary, err := s.SPM(ctx, f)
	if err != nil {
		return nil, err
	}
	if format == "" {
		format = "xlsx"
	}
	stamp := time.Now().Format("20060102-1504")
	switch strings.ToLower(format) {
	case "xlsx":
		data, err := exportXLSX(items, summary)
		if err != nil {
			return nil, err
		}
		return &ExportFile{
			Data: data, ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Filename: "mata-sidik-spm-" + stamp + ".xlsx",
		}, nil
	case "pdf":
		data, err := exportPDF(items, summary, f)
		if err != nil {
			return nil, err
		}
		return &ExportFile{Data: data, ContentType: "application/pdf", Filename: "mata-sidik-spm-" + stamp + ".pdf"}, nil
	default:
		return nil, domain.Validation("format export harus xlsx atau pdf", nil)
	}
}

func statusText(status string) string {
	switch status {
	case "ON_TARGET":
		return "Sesuai target"
	case "BELOW_TARGET":
		return "Di bawah target"
	case "NO_DATA":
		return "Belum ada data"
	case "STALE":
		return "Perlu hitung ulang"
	default:
		return status
	}
}

func exportXLSX(items []domain.DashboardSPMItem, summary *domain.DashboardSummary) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "SPM"
	_ = f.SetSheetName("Sheet1", sheet)
	baselineYear := 0
	if len(items) > 0 {
		baselineYear = items[0].BaselineYear
	}
	baselineHeader := "Capaian sebelumnya"
	if baselineYear > 0 {
		baselineHeader = fmt.Sprintf("Capaian %d", baselineYear)
	}
	headers := []string{"No", "Indikator", "Label", "Target", baselineHeader, "Capaian TW1", "Capaian", "Status", "Gap", "Num/Den", "Source hash", "Data version"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FFDB33"}},
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		Border:    []excelize.Border{{Type: "left", Color: "111111", Style: 1}, {Type: "top", Color: "111111", Style: 1}, {Type: "right", Color: "111111", Style: 1}, {Type: "bottom", Color: "111111", Style: 1}},
	})
	_ = f.SetCellStyle(sheet, "A1", "L1", headerStyle)
	for i, item := range items {
		row := i + 2
		values := []any{
			item.Kode, item.Indikator, deref(item.LabelKode), item.Target,
			formatNumber(item.BaselinePrev, 2), formatNumber(item.BaselineTW1, 2),
			formatNumber(item.Value, 2), statusText(item.DisplayStatus), formatNumber(item.Gap, 2),
			fmt.Sprintf("%s/%s", formatNumber(item.Numerator, 0), formatNumber(item.Denominator, 0)),
			deref(item.SourceHash), deref(item.DataVersion),
		}
		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			_ = f.SetCellValue(sheet, cell, v)
		}
	}
	_ = f.SetColWidth(sheet, "A", "A", 6)
	_ = f.SetColWidth(sheet, "B", "B", 38)
	_ = f.SetColWidth(sheet, "C", "C", 22)
	_ = f.SetColWidth(sheet, "D", "L", 14)
	_ = f.SetRowHeight(sheet, 1, 30)
	f.SetActiveSheet(0)
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func exportPDF(items []domain.DashboardSPMItem, summary *domain.DashboardSummary, f domain.DashboardFilter) ([]byte, error) {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(0, 8, "MATA SIDIK - Capaian SPM", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(0, 5, "Periode: "+f.PeriodeID+"  |  Total: "+fmt.Sprintf("%d", summary.Total)+"  |  Sesuai target: "+fmt.Sprintf("%d", summary.OnTarget)+"  |  Di bawah target: "+fmt.Sprintf("%d", summary.BelowTarget)+"  |  Belum ada data: "+fmt.Sprintf("%d", summary.NoData), "", 1, "L", false, 0, "")
	pdf.Ln(2)
	baselineYear := 0
	if len(items) > 0 {
		baselineYear = items[0].BaselineYear
	}
	baselineHeader := "Capaian sebelumnya"
	if baselineYear > 0 {
		baselineHeader = fmt.Sprintf("Capaian %d", baselineYear)
	}
	headers := []string{"No", "Indikator", "Target", baselineHeader, "TW1", "Capaian", "Status", "Gap", "Num/Den"}
	widths := []float64{10, 94, 22, 26, 20, 22, 36, 18, 32}
	pdf.SetFont("Helvetica", "B", 7)
	pdf.SetFillColor(255, 219, 51)
	for i, h := range headers {
		pdf.CellFormat(widths[i], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Helvetica", "", 6.8)
	for _, item := range items {
		values := []string{
			item.Kode, truncate(item.Indikator, 62), fmt.Sprintf("%.0f", item.Target),
			formatNumber(item.BaselinePrev, 1), formatNumber(item.BaselineTW1, 1),
			formatNumber(item.Value, 2), statusText(item.DisplayStatus), formatNumber(item.Gap, 2),
			fmt.Sprintf("%s/%s", formatNumber(item.Numerator, 0), formatNumber(item.Denominator, 0)),
		}
		for i, v := range values {
			align := "L"
			if i != 1 {
				align = "C"
			}
			pdf.CellFormat(widths[i], 5.2, v, "1", 0, align, false, 0, "")
		}
		pdf.Ln(-1)
	}
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatNumber(v *float64, digits int) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%.*f", digits, *v)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
