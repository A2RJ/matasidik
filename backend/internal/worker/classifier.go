package worker

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type Classifier struct {
	master repository.MasterRepository
}

func NewClassifier(master repository.MasterRepository) *Classifier {
	return &Classifier{master: master}
}

type Classified struct {
	Rows   []domain.SubmissionRowInput
	Values []domain.IndicatorValueInput
}

func (c *Classifier) Classify(ctx context.Context, sub *domain.Submission, raw []RawRow) (*Classified, error) {
	releaseID := ""
	if sub.MasterReleaseID != nil {
		releaseID = *sub.MasterReleaseID
	}
	wilayahList, err := c.master.ListWilayah(ctx)
	if err != nil {
		return nil, err
	}
	scopeByKey := map[string]domain.Wilayah{}
	for _, w := range wilayahList {
		scopeByKey[strings.ToUpper(w.Level)+"|"+w.Kode] = w
	}
	seen := map[string]bool{}
	out := &Classified{}
	for _, rr := range raw {
		row := domain.SubmissionRowInput{
			ID:               uuid.NewString(),
			RowNo:            rr.RowNo,
			KodeIKURaw:       strings.TrimSpace(rr.KodeIKU),
			KodeIndikatorRaw: strings.TrimSpace(rr.KodeIndikator),
			ScopeType:        strings.ToUpper(strings.TrimSpace(rr.ScopeType)),
			ScopeCode:        strings.TrimSpace(rr.ScopeCode),
		}
		if row.ScopeType == "" {
			row.ScopeType = "FASKES"
		}
		if row.ScopeCode == "" {
			row.ScopeCode = sub.FaskesID
		}
		if rr.Numerator != nil {
			v := *rr.Numerator
			row.Numerator = v
		}
		if rr.Denominator != nil {
			v := *rr.Denominator
			row.Denominator = v
		}
		if rr.Nilai != nil {
			v := *rr.Nilai
			row.Value = v
		}
		row.Classification = domain.ClassInvalidValue
		if row.KodeIKURaw == "" || row.KodeIndikatorRaw == "" {
			row.ErrorMessage = "kode_iku dan kode_indikator wajib diisi"
			out.Rows = append(out.Rows, row)
			continue
		}
		kodeIKUs := splitKodeIKU(row.KodeIKURaw)
		if len(kodeIKUs) == 0 {
			row.ErrorMessage = "kode_iku wajib berisi minimal satu kode"
			out.Rows = append(out.Rows, row)
			continue
		}
		if strings.TrimSpace(rr.Periode) == "" {
			row.ErrorMessage = "periode wajib diisi"
			out.Rows = append(out.Rows, row)
			continue
		}
		period, err := c.master.GetPeriodByKode(ctx, strings.TrimSpace(rr.Periode))
		if err != nil || period.ID != sub.PeriodeID {
			row.ErrorMessage = "periode baris tidak sesuai dengan submission"
			out.Rows = append(out.Rows, row)
			continue
		}
		if kodeFaskes := strings.TrimSpace(rr.KodeFaskes); kodeFaskes != "" {
			faskes, err := c.master.GetFaskes(ctx, sub.FaskesID)
			if err != nil || (kodeFaskes != sub.FaskesID && !strings.EqualFold(kodeFaskes, faskes.Kode)) {
				row.Classification = domain.ClassScopeDenied
				row.ErrorMessage = "kode_faskes tidak sesuai dengan faskes pengunggah"
				out.Rows = append(out.Rows, row)
				continue
			}
		}
		var canonicalScopeID string
		switch row.ScopeType {
		case "FASKES":
			// Legacy compatibility: old files use uploader faskes as scope.
			if row.ScopeCode != sub.FaskesID {
				row.Classification = domain.ClassScopeDenied
				row.ErrorMessage = "baris di luar scope faskes pengupload"
				out.Rows = append(out.Rows, row)
				continue
			}
			canonicalScopeID = row.ScopeCode
		case "DESA", "KECAMATAN", "KABUPATEN":
			w, ok := scopeByKey[row.ScopeType+"|"+row.ScopeCode]
			if !ok {
				row.Classification = domain.ClassScopeDenied
				row.ErrorMessage = "scope_code tidak dikenal atau levelnya tidak sesuai: " + row.ScopeCode
				out.Rows = append(out.Rows, row)
				continue
			}
			canonicalScopeID = w.ID
		default:
			row.Classification = domain.ClassScopeDenied
			row.ErrorMessage = "scope_type harus DESA, KECAMATAN, atau KABUPATEN"
			out.Rows = append(out.Rows, row)
			continue
		}
		ind, err := c.master.GetIndicatorByKode(ctx, row.KodeIndikatorRaw)
		if err != nil {
			row.Classification = domain.ClassUnknownInd
			row.ErrorMessage = "kode_indikator tidak dikenal"
			out.Rows = append(out.Rows, row)
			continue
		}
		row.MatchedIndikatorID = ind.ID
		if ind.Frekuensi == "tahunan" && period.Bulan != 0 {
			row.ErrorMessage = "indikator tahunan harus diunggah pada periode tahunan"
			out.Rows = append(out.Rows, row)
			continue
		}
		if ind.Frekuensi == "bulanan" && period.Bulan == 0 {
			row.ErrorMessage = "indikator bulanan harus diunggah pada periode bulanan"
			out.Rows = append(out.Rows, row)
			continue
		}
		tipeNilai := domain.TipeNilai(strings.ToUpper(strings.TrimSpace(ind.TipeNilai)))
		if tipeNilai == "" {
			tipeNilai = domain.TipeRatio
		}
		switch tipeNilai {
		case domain.TipeRatio:
			if rr.Numerator == nil || rr.Denominator == nil {
				row.ErrorMessage = "numerator dan denominator wajib diisi untuk indikator rasio"
				out.Rows = append(out.Rows, row)
				continue
			}
			if *rr.Numerator < 0 || *rr.Denominator < 0 {
				row.ErrorMessage = "numerator/denominator tidak boleh negatif"
				out.Rows = append(out.Rows, row)
				continue
			}
			if *rr.Denominator <= 0 {
				row.ErrorMessage = "denominator wajib lebih dari nol"
				out.Rows = append(out.Rows, row)
				continue
			}
			row.Numerator = *rr.Numerator
			row.Denominator = *rr.Denominator
			multiplier := 100.0
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(ind.Satuan)), "rasio") {
				multiplier = 1
			}
			row.Value = *rr.Numerator / *rr.Denominator * multiplier
		case domain.TipeCount, domain.TipeScalar:
			if rr.Nilai == nil {
				row.ErrorMessage = "nilai wajib diisi untuk indikator " + strings.ToLower(string(tipeNilai))
				out.Rows = append(out.Rows, row)
				continue
			}
			if *rr.Nilai < 0 {
				row.ErrorMessage = "nilai tidak boleh negatif"
				out.Rows = append(out.Rows, row)
				continue
			}
			row.Value = *rr.Nilai
			row.Numerator = 0
			row.Denominator = 0
		default:
			row.ErrorMessage = "tipe nilai indikator tidak dikenal: " + string(tipeNilai)
			out.Rows = append(out.Rows, row)
			continue
		}
		labels := make([]domain.Label, 0, len(kodeIKUs))
		for _, kodeIKU := range kodeIKUs {
			label, err := c.master.GetLabelByKode(ctx, kodeIKU, releaseID)
			if err != nil {
				row.Classification = domain.ClassUnknownIKU
				row.ErrorMessage = "kode_iku tidak dikenal pada master release: " + kodeIKU
				labels = nil
				break
			}
			if label.IsRollupOnly || !label.AllowDirectInput {
				row.Classification = domain.ClassRollupOnly
				row.ErrorMessage = "IKU induk/rollup tidak boleh menerima input langsung: " + kodeIKU
				labels = nil
				break
			}
			labels = append(labels, *label)
		}
		if len(labels) == 0 {
			out.Rows = append(out.Rows, row)
			continue
		}
		// kode_kategori, jika diisi, mengikuti label pertama (label utama/konteks).
		if rawCategory := strings.TrimSpace(rr.KodeKategori); rawCategory != "" {
			categoryKode := ""
			if labels[0].CategoryKode != nil {
				categoryKode = *labels[0].CategoryKode
			}
			if !strings.EqualFold(rawCategory, categoryKode) {
				row.Classification = domain.ClassMismatch
				row.ErrorMessage = "kode_kategori tidak sesuai dengan kategori label utama"
				out.Rows = append(out.Rows, row)
				continue
			}
		}
		links, err := c.master.ListIndicatorLabels(ctx, ind.ID)
		if err != nil {
			return nil, err
		}
		allLinked := true
		for _, label := range labels {
			linked := false
			for _, link := range links {
				if link.LabelID == label.ID {
					linked = true
					break
				}
			}
			if !linked {
				allLinked = false
				row.ErrorMessage = "salah satu kode_iku tidak terhubung ke indikator: " + label.Kode
				break
			}
		}
		if !allLinked {
			row.Classification = domain.ClassMismatch
			out.Rows = append(out.Rows, row)
			continue
		}
		// Satu nilai hanya boleh sekali per indikator + scope. Kombinasi label
		// yang berbeda tidak boleh membuat baris ganda.
		key := fmt.Sprintf("%s|%s", ind.ID, row.ScopeCode)
		if seen[key] {
			row.Classification = domain.ClassDuplicate
			row.ErrorMessage = "baris duplikat untuk indikator dan scope yang sama"
			out.Rows = append(out.Rows, row)
			continue
		}
		seen[key] = true
		row.MatchedLabelID = labels[0].ID
		row.Classification = domain.ClassMatched
		out.Rows = append(out.Rows, row)
		out.Values = append(out.Values, domain.IndicatorValueInput{
			SubmissionRowID: row.ID, IndikatorID: ind.ID, PeriodeID: sub.PeriodeID,
			FaskesID: sub.FaskesID, ScopeType: row.ScopeType, ScopeID: canonicalScopeID,
			Numerator: row.Numerator, Denominator: row.Denominator, Value: row.Value,
			Status: "COMPLETE",
		})
	}
	return out, nil
}

// splitKodeIKU memisahkan beberapa kode IKU dengan pemisah titik koma.
// Spasi di sekitar kode diabaikan; kode kosong dilewati; duplikat dipertahankan sekali.
func splitKodeIKU(raw string) []string {
	parts := strings.Split(raw, ";")
	out := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, part := range parts {
		kode := strings.TrimSpace(part)
		if kode == "" || seen[kode] {
			continue
		}
		seen[kode] = true
		out = append(out, kode)
	}
	return out
}
