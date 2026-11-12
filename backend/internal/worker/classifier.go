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
		if row.ScopeType == "FASKES" && row.ScopeCode != sub.FaskesID {
			row.Classification = domain.ClassScopeDenied
			row.ErrorMessage = "baris di luar scope faskes pengupload"
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
		if row.Numerator < 0 || row.Denominator < 0 {
			row.ErrorMessage = "numerator/denominator tidak boleh negatif"
			out.Rows = append(out.Rows, row)
			continue
		}
		if row.Denominator <= 0 {
			row.ErrorMessage = "denominator wajib lebih dari nol"
			out.Rows = append(out.Rows, row)
			continue
		}
		if row.Value == 0 {
			row.Value = row.Numerator / row.Denominator * 100
		}
		label, err := c.master.GetLabelByKode(ctx, row.KodeIKURaw, releaseID)
		if err != nil {
			row.Classification = domain.ClassUnknownIKU
			row.ErrorMessage = "kode_iku tidak dikenal pada master release"
			out.Rows = append(out.Rows, row)
			continue
		}
		if rawCategory := strings.TrimSpace(rr.KodeKategori); rawCategory != "" {
			categoryKode := ""
			if label.CategoryKode != nil {
				categoryKode = *label.CategoryKode
			}
			if !strings.EqualFold(rawCategory, categoryKode) {
				row.Classification = domain.ClassMismatch
				row.ErrorMessage = "kode_kategori tidak sesuai dengan kategori label"
				out.Rows = append(out.Rows, row)
				continue
			}
		}
		if label.IsRollupOnly || !label.AllowDirectInput {
			row.Classification = domain.ClassRollupOnly
			row.ErrorMessage = "IKU induk/rollup tidak boleh menerima input langsung"
			out.Rows = append(out.Rows, row)
			continue
		}
		links, err := c.master.ListIndicatorLabels(ctx, ind.ID)
		if err != nil {
			return nil, err
		}
		linked := false
		for _, link := range links {
			if link.LabelID == label.ID {
				linked = true
				break
			}
		}
		if !linked {
			row.Classification = domain.ClassMismatch
			row.ErrorMessage = "indikator tidak terhubung ke kode_iku tersebut"
			out.Rows = append(out.Rows, row)
			continue
		}
		key := fmt.Sprintf("%s|%s|%s", label.ID, ind.ID, row.ScopeCode)
		if seen[key] {
			row.Classification = domain.ClassDuplicate
			row.ErrorMessage = "baris duplikat dalam file"
			out.Rows = append(out.Rows, row)
			continue
		}
		seen[key] = true
		row.MatchedLabelID = label.ID
		row.Classification = domain.ClassMatched
		out.Rows = append(out.Rows, row)
		out.Values = append(out.Values, domain.IndicatorValueInput{
			SubmissionRowID: row.ID, IndikatorID: ind.ID, PeriodeID: sub.PeriodeID,
			FaskesID: sub.FaskesID, ScopeType: row.ScopeType, ScopeID: row.ScopeCode,
			Numerator: row.Numerator, Denominator: row.Denominator, Value: row.Value,
			Status: "COMPLETE",
		})
	}
	return out, nil
}
