package worker

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type RawRow struct {
	RowNo         int
	Periode       string
	KodeFaskes    string
	KodeKategori  string
	KodeIKU       string
	KodeIndikator string
	ScopeType     string
	ScopeCode     string
	Nilai         *float64
	Numerator     *float64
	Denominator   *float64
	Catatan       string
}

func ParseFile(path string) ([]RawRow, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".csv":
		return parseCSV(path)
	case ".xlsx":
		return parseXLSX(path)
	default:
		return nil, fmt.Errorf("format file tidak didukung: %s", filepath.Ext(path))
	}
}

func parseCSV(path string) ([]RawRow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	all, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	text := strings.TrimPrefix(string(all), "\ufeff")
	delim := detectDelimiter(text)
	r := csv.NewReader(strings.NewReader(text))
	r.Comma = delim
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("baca csv: %w", err)
	}
	return rowsFromRecords(records)
}

func detectDelimiter(text string) rune {
	first := text
	if i := strings.IndexByte(first, '\n'); i >= 0 {
		first = first[:i]
	}
	if strings.Count(first, ";") >= strings.Count(first, ",") {
		return ';'
	}
	return ','
}

func parseXLSX(path string) ([]RawRow, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("buka xlsx: %w", err)
	}
	defer f.Close()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("xlsx tidak punya sheet")
	}
	records, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}
	return rowsFromRecords(records)
}

func rowsFromRecords(records [][]string) ([]RawRow, error) {
	if len(records) < 2 {
		return nil, fmt.Errorf("file tidak memiliki baris data")
	}
	headerIdx := 0
	for headerIdx < len(records) && isEmptyRecord(records[headerIdx]) {
		headerIdx++
	}
	if headerIdx >= len(records) {
		return nil, fmt.Errorf("header tidak ditemukan")
	}
	headers := map[string]int{}
	for i, h := range records[headerIdx] {
		headers[normalizeHeader(h)] = i
	}
	required := []string{"kode_iku", "kode_indikator", "periode"}
	for _, key := range required {
		if _, ok := headers[key]; !ok {
			return nil, fmt.Errorf("header wajib tidak ditemukan: %s", key)
		}
	}
	var out []RawRow
	for i := headerIdx + 1; i < len(records); i++ {
		rec := records[i]
		if isEmptyRecord(rec) {
			continue
		}
		row := RawRow{
			RowNo:         len(out) + 1,
			Periode:       cell(rec, headers, "periode"),
			KodeFaskes:    cell(rec, headers, "kode_faskes", "faskes"),
			KodeKategori:  cell(rec, headers, "kode_kategori", "kategori"),
			KodeIKU:       cell(rec, headers, "kode_iku", "iku"),
			KodeIndikator: cell(rec, headers, "kode_indikator", "indikator"),
			ScopeType:     strings.ToUpper(cell(rec, headers, "scope_type")),
			ScopeCode:     cell(rec, headers, "scope_code"),
			Catatan:       cell(rec, headers, "catatan", "notes"),
		}
		row.Nilai = numberCell(rec, headers, "nilai", "value")
		row.Numerator = numberCell(rec, headers, "numerator", "num", "pembilang")
		row.Denominator = numberCell(rec, headers, "denominator", "den", "penyebut")
		out = append(out, row)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("file tidak memiliki baris data")
	}
	return out, nil
}

func normalizeHeader(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

func cell(rec []string, headers map[string]int, keys ...string) string {
	for _, key := range keys {
		if idx, ok := headers[key]; ok && idx < len(rec) {
			return strings.TrimSpace(rec[idx])
		}
	}
	return ""
}

func numberCell(rec []string, headers map[string]int, keys ...string) *float64 {
	v := cell(rec, headers, keys...)
	if v == "" {
		return nil
	}
	f, err := parseNumber(v)
	if err != nil {
		return nil
	}
	return &f
}

func parseNumber(v string) (float64, error) {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, "%", "")
	v = strings.ReplaceAll(v, " ", "")
	if strings.Contains(v, ",") && strings.Contains(v, ".") {
		v = strings.ReplaceAll(v, ",", "")
	} else if strings.Contains(v, ",") {
		v = strings.ReplaceAll(v, ",", ".")
	}
	return strconv.ParseFloat(v, 64)
}

func isEmptyRecord(rec []string) bool {
	for _, c := range rec {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}
