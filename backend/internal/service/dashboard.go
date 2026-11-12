package service

import (
	"context"
	"strings"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type DashboardService struct {
	master repository.MasterRepository
	values repository.ValueRepository
	subs   repository.SubmissionRepository
}

func NewDashboardService(master repository.MasterRepository, values repository.ValueRepository, subs repository.SubmissionRepository) *DashboardService {
	return &DashboardService{master: master, values: values, subs: subs}
}

func (s *DashboardService) ResolveFilter(ctx context.Context, f domain.DashboardFilter) (domain.DashboardFilter, error) {
	if strings.TrimSpace(f.PeriodeID) == "" {
		if id, err := s.master.DefaultPeriodID(ctx); err == nil && id != "" {
			f.PeriodeID = id
		} else {
			periods, err := s.master.ListPeriods(ctx)
			if err != nil {
				return f, err
			}
			if len(periods) > 0 {
				f.PeriodeID = periods[0].ID
			}
		}
	} else if len(f.PeriodeID) == 7 && strings.Count(f.PeriodeID, "-") == 1 {
		if p, err := s.master.GetPeriodByKode(ctx, f.PeriodeID); err == nil {
			f.PeriodeID = p.ID
		}
	}
	if f.FaskesID != "" {
		f.ScopeType = "FASKES"
		f.ScopeID = f.FaskesID
	}
	if f.ScopeType == "" {
		f.ScopeType = "KABUPATEN"
	}
	if f.ScopeID == "" && f.ScopeType == "KABUPATEN" {
		wilayah, err := s.master.ListWilayah(ctx)
		if err == nil {
			for _, w := range wilayah {
				if w.Level == "kabupaten" {
					f.ScopeID = w.ID
					break
				}
			}
		}
	}
	return f, nil
}

func (s *DashboardService) SPM(ctx context.Context, f domain.DashboardFilter) ([]domain.DashboardSPMItem, *domain.DashboardSummary, error) {
	f, err := s.ResolveFilter(ctx, f)
	if err != nil {
		return nil, nil, err
	}
	items, err := s.values.DashboardSPM(ctx, f)
	if err != nil {
		return nil, nil, err
	}
	labels, _ := s.master.ListLabels(ctx)
	labelByID := map[string]domain.Label{}
	for _, l := range labels {
		labelByID[l.ID] = l
	}
	for i := range items {
		item := &items[i]
		filteredOut := false
		if item.LabelID != nil {
			if l, ok := labelByID[*item.LabelID]; ok {
				item.LabelKode = &l.Kode
				item.LabelNama = &l.Nama
				item.CategoryKode = l.CategoryKode
				item.CategoryNama = l.CategoryNama
				if f.BidangID != "" && (l.PJBidangID == nil || *l.PJBidangID != f.BidangID) {
					filteredOut = true
				}
			}
		}
		if item.Value != nil && item.Target > 0 {
			gap := *item.Value - item.Target
			item.Gap = &gap
		}
		item.DisplayStatus = displayStatus(*item)
		if filteredOut {
			item.DisplayStatus = "FILTERED"
		}
	}
	filtered := items[:0]
	for _, item := range items {
		if item.DisplayStatus == "FILTERED" {
			continue
		}
		if f.IndikatorID != "" && (item.IndikatorID == nil || *item.IndikatorID != f.IndikatorID) {
			continue
		}
		if f.Status != "" && item.DisplayStatus != f.Status {
			continue
		}
		filtered = append(filtered, item)
	}
	summary := summarize(filtered)
	// Nilai IKU induk menjadi kartu ringkasan utama bila tersedia.
	for _, l := range labels {
		if l.Kode == "IKU-2026-001" {
			if v, _ := s.values.GetIKUValue(ctx, l.ID, f.PeriodeID, f.ScopeType, f.ScopeID); v != nil {
				summary.ParentValue = v.Value
				summary.ParentStatus = string(v.Status)
			}
			break
		}
	}
	return filtered, summary, nil
}

func (s *DashboardService) Category(ctx context.Context, kode string, f domain.DashboardFilter) (*domain.DashboardCategoryResult, error) {
	f, err := s.ResolveFilter(ctx, f)
	if err != nil {
		return nil, err
	}
	kode = strings.ToUpper(strings.TrimSpace(kode))
	if kode == "ALL" {
		return s.allCategories(ctx, f)
	}
	category, err := s.master.GetDashboardCategoryByKode(ctx, kode)
	if err != nil {
		return nil, err
	}
	result := &domain.DashboardCategoryResult{Category: *category, Items: []domain.DashboardSPMItem{}}
	if category.Kode == "SPM" {
		items, summary, err := s.SPM(ctx, f)
		if err != nil {
			return nil, err
		}
		result.Items = items
		result.Summary = summary
		return result, nil
	}
	targets, err := s.master.ListTargets(ctx, f.PeriodeID)
	if err != nil {
		return nil, err
	}
	labelTargets, indicatorTargets := targetMaps(targets, f)
	if category.Jenis == "OWNER" {
		return s.ownerCategory(ctx, category, f, indicatorTargets, result)
	}
	if err := s.appendLabelCategory(ctx, category, f, labelTargets, result); err != nil {
		return nil, err
	}
	result.Summary = summarize(result.Items)
	return result, nil
}

func (s *DashboardService) allCategories(ctx context.Context, f domain.DashboardFilter) (*domain.DashboardCategoryResult, error) {
	categories, err := s.master.ListDashboardCategories(ctx)
	if err != nil {
		return nil, err
	}
	targets, err := s.master.ListTargets(ctx, f.PeriodeID)
	if err != nil {
		return nil, err
	}
	labelTargets, indicatorTargets := targetMaps(targets, f)
	result := &domain.DashboardCategoryResult{
		Category: domain.DashboardCategory{
			ID: "all", Kode: "ALL", Nama: "Semua", Jenis: "LABEL", Urutan: 0, Aktif: true, IsSystem: true,
		},
		Items: []domain.DashboardSPMItem{},
	}
	for _, category := range categories {
		if !category.Aktif {
			continue
		}
		category := category
		if category.Jenis == "OWNER" {
			if _, err := s.ownerCategory(ctx, &category, f, indicatorTargets, result); err != nil {
				return nil, err
			}
			continue
		}
		if err := s.appendLabelCategory(ctx, &category, f, labelTargets, result); err != nil {
			return nil, err
		}
	}
	result.Summary = summarize(result.Items)
	return result, nil
}

func (s *DashboardService) appendLabelCategory(ctx context.Context, category *domain.DashboardCategory, f domain.DashboardFilter, labelTargets map[string]float64, result *domain.DashboardCategoryResult) error {
	labels, err := s.master.ListLabels(ctx)
	if err != nil {
		return err
	}
	for _, label := range labels {
		if label.CategoryID == nil || *label.CategoryID != category.ID {
			continue
		}
		item := domain.DashboardSPMItem{
			Kode: label.Kode, Indikator: label.Nama, LabelID: &label.ID,
			LabelKode: &label.Kode, LabelNama: &label.Nama,
			CategoryKode: &category.Kode, CategoryNama: &category.Nama,
			Target: labelTargets[label.ID], ValueStatus: string(domain.ValueNoData),
		}
		if v, err := s.values.GetIKUValue(ctx, label.ID, f.PeriodeID, f.ScopeType, f.ScopeID); err == nil && v != nil {
			item.IKUValueID = &v.ID
			item.Value = v.Value
			item.Numerator = &v.Numerator
			item.Denominator = &v.Denominator
			item.ValueStatus = string(v.Status)
			item.SourceHash = v.SourceHash
			item.DataVersion = v.DataVersion
			item.CalculationRunID = v.CalculationRunID
			item.UpdatedAt = &v.UpdatedAt
		}
		item.DisplayStatus = displayStatus(item)
		if item.Value != nil && item.Target > 0 {
			gap := *item.Value - item.Target
			item.Gap = &gap
		}
		if f.Status != "" && item.DisplayStatus != f.Status {
			continue
		}
		result.Items = append(result.Items, item)
	}
	return nil
}

func (s *DashboardService) ownerCategory(ctx context.Context, category *domain.DashboardCategory, f domain.DashboardFilter, indicatorTargets map[string]float64, result *domain.DashboardCategoryResult) (*domain.DashboardCategoryResult, error) {
	indicators, err := s.master.ListIndicators(ctx)
	if err != nil {
		return nil, err
	}
	bidangList, err := s.master.ListBidang(ctx)
	if err != nil {
		return nil, err
	}
	bidangByID := map[string]domain.Bidang{}
	for _, b := range bidangList {
		bidangByID[b.ID] = b
	}
	values, err := s.values.ListFinalIndicatorValuesForScope(ctx, f.PeriodeID, f.ScopeType, f.ScopeID)
	if err != nil {
		return nil, err
	}
	type aggregate struct {
		num float64
		den float64
	}
	byIndicator := map[string]aggregate{}
	for _, v := range values {
		a := byIndicator[v.IndikatorID]
		a.num += v.Numerator
		a.den += v.Denominator
		byIndicator[v.IndikatorID] = a
	}
	for _, indicator := range indicators {
		if f.TimkerID != "" && (indicator.PJTimkerID == nil || *indicator.PJTimkerID != f.TimkerID) {
			continue
		}
		if f.BidangID != "" && !indicatorOwnedByBidang(indicator, f.BidangID, bidangByID) {
			continue
		}
		item := domain.DashboardSPMItem{
			Kode: indicator.Kode, Indikator: indicator.Nama, IndikatorID: &indicator.ID,
			CategoryKode: &category.Kode, CategoryNama: &category.Nama,
			Target: indicatorTargets[indicator.ID], ValueStatus: string(domain.ValueNoData),
		}
		if a, ok := byIndicator[indicator.ID]; ok && a.den > 0 {
			value := a.num / a.den * 100
			item.Value = &value
			item.Numerator = &a.num
			item.Denominator = &a.den
			item.ValueStatus = string(domain.ValueComplete)
		}
		ownerID := indicator.PJBidangID
		if indicator.PJTimkerID != nil {
			ownerID = indicator.PJTimkerID
		}
		if ownerID != nil {
			if owner, ok := bidangByID[*ownerID]; ok {
				item.OwnerID = &owner.ID
				item.OwnerKode = &owner.Kode
				item.OwnerNama = &owner.Nama
				item.OwnerJenis = &owner.Jenis
			}
		}
		item.DisplayStatus = displayStatus(item)
		if item.Value != nil && item.Target > 0 {
			gap := *item.Value - item.Target
			item.Gap = &gap
		}
		if f.Status != "" && item.DisplayStatus != f.Status {
			continue
		}
		result.Items = append(result.Items, item)
	}
	result.Summary = summarize(result.Items)
	return result, nil
}

func indicatorOwnedByBidang(indicator domain.Indikator, bidangID string, bidangByID map[string]domain.Bidang) bool {
	if indicator.PJBidangID != nil && *indicator.PJBidangID == bidangID {
		return true
	}
	if indicator.PJTimkerID == nil {
		return false
	}
	timker, ok := bidangByID[*indicator.PJTimkerID]
	return ok && timker.ParentID != nil && *timker.ParentID == bidangID
}

func targetMaps(targets []domain.Target, f domain.DashboardFilter) (map[string]float64, map[string]float64) {
	type picked struct {
		value float64
		exact bool
	}
	labelPick := map[string]picked{}
	indicatorPick := map[string]picked{}
	set := func(m map[string]picked, id string, target float64, exact bool) {
		if id == "" {
			return
		}
		current, ok := m[id]
		if !ok || (exact && !current.exact) {
			m[id] = picked{value: target, exact: exact}
		}
	}
	for _, t := range targets {
		exact := t.ScopeType == f.ScopeType && deref(t.ScopeID) == f.ScopeID
		if !exact && t.ScopeType != "KABUPATEN" {
			continue
		}
		if t.LabelID != nil {
			set(labelPick, *t.LabelID, t.TargetBulanan, exact)
		}
		if t.IndikatorID != nil {
			set(indicatorPick, *t.IndikatorID, t.TargetBulanan, exact)
		}
	}
	labels := map[string]float64{}
	for id, p := range labelPick {
		labels[id] = p.value
	}
	indicators := map[string]float64{}
	for id, p := range indicatorPick {
		indicators[id] = p.value
	}
	return labels, indicators
}

func (s *DashboardService) Trace(ctx context.Context, kode string, f domain.DashboardFilter) (*domain.SPMTrace, error) {
	f, err := s.ResolveFilter(ctx, f)
	if err != nil {
		return nil, err
	}
	trace, err := s.values.TraceSPM(ctx, kode, f)
	if err != nil {
		return nil, err
	}
	trace.Row.DisplayStatus = displayStatus(trace.Row)
	if trace.Row.Value != nil && trace.Row.Target > 0 {
		gap := *trace.Row.Value - trace.Row.Target
		trace.Row.Gap = &gap
	}
	return trace, nil
}

func (s *DashboardService) IKUValues(ctx context.Context, periodeID, scopeType, scopeID string) ([]domain.IKUValue, error) {
	if periodeID == "" {
		periods, _ := s.master.ListPeriods(ctx)
		if len(periods) > 0 {
			periodeID = periods[0].ID
		}
	}
	if scopeType == "" {
		scopeType = "KABUPATEN"
	}
	if scopeID == "" {
		scopeID = "wil-kab-swb"
	}
	return s.values.ListIKUValues(ctx, periodeID, scopeType, scopeID)
}

func displayStatus(item domain.DashboardSPMItem) string {
	if item.Value == nil || item.ValueStatus == string(domain.ValueNoData) || item.ValueStatus == "" {
		return "NO_DATA"
	}
	if item.ValueStatus == string(domain.ValueStale) {
		return "STALE"
	}
	if item.Target <= 0 {
		return "NO_TARGET"
	}
	if *item.Value >= item.Target {
		return "ON_TARGET"
	}
	return "BELOW_TARGET"
}

func summarize(items []domain.DashboardSPMItem) *domain.DashboardSummary {
	out := &domain.DashboardSummary{Total: len(items)}
	sum := 0.0
	count := 0
	for _, item := range items {
		switch item.DisplayStatus {
		case "ON_TARGET":
			out.OnTarget++
		case "BELOW_TARGET":
			out.BelowTarget++
		case "NO_TARGET":
			out.NoTarget++
		default:
			out.NoData++
		}
		if item.Value != nil {
			sum += *item.Value
			count++
		}
	}
	if count > 0 {
		out.Average = sum / float64(count)
	}
	return out
}

func (s *DashboardService) Detail(ctx context.Context, kode string, f domain.DashboardFilter) (*domain.SPMDetail, error) {
	f, err := s.ResolveFilter(ctx, f)
	if err != nil {
		return nil, err
	}
	rowDef, err := s.master.GetDashboardRowByKode(ctx, kode)
	if err != nil {
		return nil, err
	}
	items, err := s.values.DashboardSPM(ctx, f)
	if err != nil {
		return nil, err
	}
	var row *domain.DashboardSPMItem
	for i := range items {
		if items[i].Kode == kode {
			row = &items[i]
			break
		}
	}
	if row == nil {
		return nil, domain.NotFound("baris dashboard")
	}
	row.DisplayStatus = displayStatus(*row)
	if row.Value != nil && row.Target > 0 {
		gap := *row.Value - row.Target
		row.Gap = &gap
	}
	detail := &domain.SPMDetail{Row: *row}
	if row.IKUValueID != nil {
		provenance, err := s.values.ListIKUValueDetails(ctx, *row.IKUValueID)
		if err == nil {
			detail.Provenance = provenance
		}
	}
	faskesList, err := s.master.ListFaskes(ctx)
	if err != nil {
		return nil, err
	}
	for _, fk := range faskesList {
		item := domain.SPMFaskesDetail{
			FaskesID: fk.ID, FaskesKode: fk.Kode, FaskesNama: fk.Nama, Status: string(domain.ValueNoData),
		}
		if rowDef.LabelID != nil {
			if v, err := s.values.GetIKUValue(ctx, *rowDef.LabelID, f.PeriodeID, "FASKES", fk.ID); err == nil && v != nil {
				item.Numerator = &v.Numerator
				item.Denominator = &v.Denominator
				item.Value = v.Value
				item.Status = string(v.Status)
				item.UpdatedAt = &v.UpdatedAt
				item.Uploaded = v.Value != nil && v.Status != domain.ValueNoData
			}
		}
		if rowDef.IndikatorID != nil {
			vals, _ := s.values.ListFinalIndicatorValuesForScope(ctx, f.PeriodeID, "FASKES", fk.ID)
			for _, iv := range vals {
				if iv.IndikatorID != *rowDef.IndikatorID {
					continue
				}
				item.Uploaded = true
				sr, err := s.subs.GetRow(ctx, iv.SubmissionRowID)
				if err != nil {
					continue
				}
				sub, err := s.subs.GetSubmission(ctx, sr.SubmissionID)
				if err != nil {
					continue
				}
				item.SubmissionID = &sub.ID
				item.SubmissionKode = &sub.Kode
				item.UploaderNama = sub.UploaderNama
				file, err := s.subs.GetFileObject(ctx, sub.FileID)
				if err == nil {
					item.FileID = &file.ID
					item.FileName = &file.OriginalName
				}
				break
			}
		}
		detail.Faskes = append(detail.Faskes, item)
	}
	return detail, nil
}

func (s *DashboardService) IKUByFaskes(ctx context.Context, f domain.DashboardFilter) (*domain.IKUChart, error) {
	f, err := s.ResolveFilter(ctx, f)
	if err != nil {
		return nil, err
	}
	faskesList, err := s.master.ListFaskes(ctx)
	if err != nil {
		return nil, err
	}
	labels, err := s.master.ListLabels(ctx)
	if err != nil {
		return nil, err
	}
	var ikuLabels []domain.Label
	for _, l := range labels {
		if l.Tipe == "IKU" && l.AllowDirectInput && !l.IsRollupOnly {
			ikuLabels = append(ikuLabels, l)
		}
	}
	chart := &domain.IKUChart{}
	for _, fk := range faskesList {
		chart.Faskes = append(chart.Faskes, domain.ChartFaskes{ID: fk.ID, Kode: fk.Kode, Nama: fk.Nama})
	}
	for _, label := range ikuLabels {
		series := domain.ChartSeries{Kode: label.Kode, Nama: label.Nama}
		for _, fk := range faskesList {
			var val *float64
			if v, err := s.values.GetIKUValue(ctx, label.ID, f.PeriodeID, "FASKES", fk.ID); err == nil && v != nil {
				val = v.Value
			}
			series.Values = append(series.Values, val)
		}
		chart.Series = append(chart.Series, series)
	}
	return chart, nil
}
