package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type CalculationService struct {
	master repository.MasterRepository
	values repository.ValueRepository
}

func NewCalculationService(master repository.MasterRepository, values repository.ValueRepository) *CalculationService {
	return &CalculationService{master: master, values: values}
}

type labelComputation struct {
	status  domain.ValueStatus
	kind    domain.TipeNilai
	scale   float64
	num     float64
	den     float64
	value   *float64
	details []domain.IKUValueDetail
}

func (s *CalculationService) CalculatePeriod(ctx context.Context, periodeID, trigger string) ([]domain.CalculationRun, error) {
	period, err := s.master.GetPeriodByID(ctx, periodeID)
	if err != nil {
		return nil, err
	}
	wilayah, err := s.master.ListWilayah(ctx)
	if err != nil {
		return nil, err
	}
	kabupatenID := ""
	for _, w := range wilayah {
		if w.Level == "kabupaten" {
			kabupatenID = w.ID
			break
		}
	}
	type scope struct {
		Type string
		ID   string
	}
	scopeSet := map[string]scope{}
	addScope := func(t, id string) {
		if t != "" && id != "" {
			scopeSet[t+"|"+id] = scope{Type: t, ID: id}
		}
	}
	addScope("KABUPATEN", kabupatenID)

	allValues, err := s.values.ListFinalIndicatorValues(ctx, periodeID)
	if err != nil {
		return nil, err
	}
	mappings, err := s.master.ListDesaFaskesMappings(ctx)
	if err != nil {
		return nil, err
	}
	mappedByFaskes := map[string]map[string]bool{}
	for _, m := range mappings {
		if !m.Aktif {
			continue
		}
		if mappedByFaskes[m.FaskesID] == nil {
			mappedByFaskes[m.FaskesID] = map[string]bool{}
		}
		mappedByFaskes[m.FaskesID][m.DesaID] = true
	}
	wilayahByID := map[string]domain.Wilayah{}
	for _, w := range wilayah {
		wilayahByID[w.ID] = w
	}
	for _, m := range mappings {
		if m.Aktif {
			addScope("FASKES", m.FaskesID)
		}
	}
	for _, v := range allValues {
		switch v.ScopeType {
		case "DESA":
			if v.ScopeID != nil {
				addScope("DESA", *v.ScopeID)
				if w, ok := wilayahByID[*v.ScopeID]; ok && w.ParentID != nil {
					addScope("KECAMATAN", *w.ParentID)
				}
			}
		case "KECAMATAN":
			if v.ScopeID != nil {
				addScope("KECAMATAN", *v.ScopeID)
				if w, ok := wilayahByID[*v.ScopeID]; ok && w.ParentID != nil {
					addScope("KABUPATEN", *w.ParentID)
				}
			}
		case "KABUPATEN":
			if v.ScopeID != nil {
				addScope("KABUPATEN", *v.ScopeID)
			}
		}
		if v.FaskesID != nil && *v.FaskesID != "" {
			addScope("FASKES", *v.FaskesID)
		}
	}
	scopes := make([]scope, 0, len(scopeSet))
	// KABUPATEN first, then geographic, then FASKES; deterministic.
	for _, sc := range scopeSet {
		if sc.Type == "KABUPATEN" {
			scopes = append(scopes, sc)
		}
	}
	for _, sc := range scopeSet {
		if sc.Type != "KABUPATEN" && sc.Type != "FASKES" {
			scopes = append(scopes, sc)
		}
	}
	for _, sc := range scopeSet {
		if sc.Type == "FASKES" {
			scopes = append(scopes, sc)
		}
	}
	sort.SliceStable(scopes, func(i, j int) bool {
		if scopes[i].Type == scopes[j].Type {
			return scopes[i].ID < scopes[j].ID
		}
		return scopes[i].Type < scopes[j].Type
	})

	var runs []domain.CalculationRun
	for _, sc := range scopes {
		values := selectScopeValues(allValues, wilayah, mappedByFaskes, sc.Type, sc.ID)
		run, err := s.calculateScopeWithValues(ctx, period, sc.Type, sc.ID, trigger, values)
		if err != nil {
			return runs, err
		}
		runs = append(runs, *run)
	}
	return runs, nil
}

// selectScopeValues chooses values for a scope. Geographic scopes use the
// lowest available level per indicator. KABUPATEN falls back to legacy FASKES
// values only when there is no geographic value for that indicator.
func selectScopeValues(all []domain.IndicatorValue, wilayah []domain.Wilayah, mappedByFaskes map[string]map[string]bool, scopeType, scopeID string) []domain.IndicatorValue {
	byID := map[string]domain.Wilayah{}
	for _, w := range wilayah {
		byID[w.ID] = w
	}
	descendants := map[string]bool{}
	if scopeType == "FASKES" {
		for _, v := range all {
			if v.FaskesID != nil && *v.FaskesID == scopeID {
				descendants["*"] = true
			}
		}
	} else {
		switch scopeType {
		case "KABUPATEN":
			descendants["*"] = true
		case "KECAMATAN":
			descendants[scopeID] = true
			for _, w := range wilayah {
				if w.ParentID != nil && *w.ParentID == scopeID {
					descendants[w.ID] = true
				}
			}
		case "DESA":
			descendants[scopeID] = true
		}
	}

	geo := map[string][]domain.IndicatorValue{}
	faskes := map[string][]domain.IndicatorValue{}
	for _, v := range all {
		if v.FaskesID != nil && *v.FaskesID != "" {
			faskes[v.IndikatorID] = append(faskes[v.IndikatorID], v)
		}
		if scopeType == "FASKES" {
			continue
		}
		if v.ScopeType == "FASKES" {
			continue
		}
		if v.ScopeID == nil {
			continue
		}
		if descendants["*"] || descendants[*v.ScopeID] {
			geo[v.IndikatorID] = append(geo[v.IndikatorID], v)
		}
	}

	out := make([]domain.IndicatorValue, 0)
	if scopeType == "FASKES" {
		mapped := mappedByFaskes[scopeID]
		geoByIndicator := map[string][]domain.IndicatorValue{}
		directByIndicator := map[string][]domain.IndicatorValue{}
		for _, v := range all {
			if v.ScopeType == "DESA" && v.ScopeID != nil && mapped[*v.ScopeID] {
				geoByIndicator[v.IndikatorID] = append(geoByIndicator[v.IndikatorID], v)
			}
			if v.ScopeType == "FASKES" && v.FaskesID != nil && *v.FaskesID == scopeID {
				directByIndicator[v.IndikatorID] = append(directByIndicator[v.IndikatorID], v)
			}
		}
		for indID, vals := range geoByIndicator {
			out = append(out, vals...)
			delete(directByIndicator, indID)
		}
		for _, vals := range directByIndicator {
			out = append(out, vals...)
		}
		return out
	}

	// For each indicator, use the lowest geographic level present.
	for _, vals := range geo {
		best := 99
		for _, v := range vals {
			if p := scopePriority(v.ScopeType); p < best {
				best = p
			}
		}
		for _, v := range vals {
			if scopePriority(v.ScopeType) == best {
				out = append(out, v)
			}
		}
	}
	if scopeType == "KABUPATEN" {
		// Legacy seed/import files store FASKES values. Use them only for
		// indicators that have no geographic values at all.
		for indID, vals := range faskes {
			if len(geo[indID]) > 0 {
				continue
			}
			out = append(out, vals...)
		}
	}
	return out
}

func scopePriority(scopeType string) int {
	switch scopeType {
	case "DESA":
		return 0
	case "KECAMATAN":
		return 1
	case "KABUPATEN":
		return 2
	case "FASKES":
		return 3
	default:
		return 4
	}
}

func (s *CalculationService) calculateScopeWithValues(ctx context.Context, period *domain.Periode, scopeType, scopeID, trigger string, values []domain.IndicatorValue) (*domain.CalculationRun, error) {
	run := &domain.CalculationRun{
		PeriodeID: period.ID, ScopeType: scopeType, ScopeID: strPtr(scopeID),
		Trigger: trigger, Status: "running", StartedAt: now(),
	}
	if err := s.values.CreateCalculationRun(ctx, run); err != nil {
		return nil, err
	}
	labels, err := s.master.ListLabels(ctx)
	if err != nil {
		s.values.FinishCalculationRun(ctx, run.ID, "failed", err.Error(), now())
		return nil, err
	}
	labelMap := make(map[string]domain.Label, len(labels))
	for _, l := range labels {
		labelMap[l.ID] = l
	}
	indicatorValues := map[string][]domain.IndicatorValue{}
	for _, v := range values {
		indicatorValues[v.IndikatorID] = append(indicatorValues[v.IndikatorID], v)
	}
	children := map[string][]domain.IKUComponent{}
	for _, l := range labels {
		comps, err := s.master.ListIKUComponents(ctx, l.ID)
		if err != nil {
			s.values.FinishCalculationRun(ctx, run.ID, "failed", err.Error(), now())
			return nil, err
		}
		children[l.ID] = comps
	}
	indicatorLinks := map[string][]string{}
	for _, l := range labels {
		links, err := s.master.ListIndicatorLabelsForLabel(ctx, l.ID)
		if err != nil {
			s.values.FinishCalculationRun(ctx, run.ID, "failed", err.Error(), now())
			return nil, err
		}
		ids := make([]string, 0, len(links))
		for _, link := range links {
			ids = append(ids, link.IndikatorID)
		}
		indicatorLinks[l.ID] = ids
	}
	indicatorList, err := s.master.ListIndicators(ctx)
	if err != nil {
		s.values.FinishCalculationRun(ctx, run.ID, "failed", err.Error(), now())
		return nil, err
	}
	indicatorByID := make(map[string]domain.Indikator, len(indicatorList))
	for _, ind := range indicatorList {
		indicatorByID[ind.ID] = ind
	}
	memo := map[string]labelComputation{}
	visiting := map[string]bool{}
	var compute func(labelID string) labelComputation
	compute = func(labelID string) labelComputation {
		if c, ok := memo[labelID]; ok {
			return c
		}
		if visiting[labelID] {
			return labelComputation{status: domain.ValueNoData}
		}
		visiting[labelID] = true
		defer delete(visiting, labelID)
		label := labelMap[labelID]
		var c labelComputation
		if label.AllowDirectInput && !label.IsRollupOnly {
			var (
				ratioNum, ratioDen            float64
				ratioScale                    float64
				sum                           float64
				hasRatio, hasCount, hasScalar bool
				latest                        *domain.IndicatorValue
			)
			for _, indicatorID := range indicatorLinks[labelID] {
				ind := indicatorByID[indicatorID]
				tipe := domain.TipeNilai(strings.ToUpper(strings.TrimSpace(ind.TipeNilai)))
				if tipe == "" {
					tipe = domain.TipeRatio
				}
				for _, iv := range indicatorValues[indicatorID] {
					switch tipe {
					case domain.TipeRatio:
						ratioNum += iv.Numerator
						ratioDen += iv.Denominator
						hasRatio = true
						scale := ratioMultiplier(ind.Satuan)
						if ratioScale == 0 {
							ratioScale = scale
						} else if ratioScale != scale {
							ratioScale = 100
						}
						componentValue := 0.0
						if iv.Denominator != 0 {
							componentValue = iv.Numerator / iv.Denominator * scale
						}
						c.details = append(c.details, domain.IKUValueDetail{
							ComponentType: "INDIKATOR", ComponentID: iv.IndikatorID,
							ComponentValue: &componentValue, Numerator: &iv.Numerator, Denominator: &iv.Denominator,
							Weight: 1, Contribution: &componentValue,
						})
					case domain.TipeCount:
						sum += iv.Value
						hasCount = true
						componentValue := iv.Value
						zero := 0.0
						c.details = append(c.details, domain.IKUValueDetail{
							ComponentType: "INDIKATOR", ComponentID: iv.IndikatorID,
							ComponentValue: &componentValue, Numerator: &componentValue, Denominator: &zero,
							Weight: 1, Contribution: &componentValue,
						})
					case domain.TipeScalar:
						hasScalar = true
						if latest == nil || iv.CreatedAt > latest.CreatedAt {
							copy := iv
							latest = &copy
						}
					}
				}
			}
			switch {
			case hasRatio && ratioDen > 0:
				c.kind = domain.TipeRatio
				if ratioScale == 0 {
					ratioScale = 100
				}
				c.scale = ratioScale
				c.num, c.den = ratioNum, ratioDen
				v := ratioNum / ratioDen * ratioScale
				c.value, c.status = &v, domain.ValueComplete
			case hasCount:
				c.kind = domain.TipeCount
				c.scale = 1
				c.num, c.den = sum, 0
				v := sum
				c.value, c.status = &v, domain.ValueComplete
			case hasScalar && latest != nil:
				c.kind = domain.TipeScalar
				c.scale = 1
				c.num, c.den = latest.Value, 0
				v := latest.Value
				c.value, c.status = &v, domain.ValueComplete
				zero := 0.0
				c.details = append(c.details, domain.IKUValueDetail{
					ComponentType: "INDIKATOR", ComponentID: latest.IndikatorID,
					ComponentValue: &v, Numerator: &v, Denominator: &zero,
					Weight: 1, Contribution: &v,
				})
			default:
				c.status = domain.ValueNoData
			}
		} else if len(children[labelID]) > 0 {
			allComplete := true
			any := false
			sum := 0.0
			ratioNum, ratioDen := 0.0, 0.0
			ratioScale := 0.0
			kinds := map[domain.TipeNilai]bool{}
			for _, comp := range children[labelID] {
				childID := comp.KomponenID
				child := compute(childID)
				if child.status == domain.ValueComplete {
					any = true
				} else {
					allComplete = false
				}
				if child.value != nil {
					sum += *child.value
				}
				if child.kind == domain.TipeRatio {
					ratioNum += child.num
					ratioDen += child.den
					if child.scale != 0 {
						if ratioScale == 0 {
							ratioScale = child.scale
						} else if ratioScale != child.scale {
							ratioScale = 100
						}
					}
				}
				kinds[child.kind] = true
				var cv *float64
				if child.value != nil {
					cv = child.value
				}
				contribution := 0.0
				if cv != nil {
					contribution = *cv * comp.Bobot
				}
				c.details = append(c.details, domain.IKUValueDetail{
					ComponentType: "CHILD_IKU", ComponentID: childID, ComponentLabelID: &childID,
					ComponentValue: cv, Numerator: &child.num, Denominator: &child.den,
					Weight: comp.Bobot, Contribution: &contribution,
				})
			}
			useSum := label.AggregationMethod == domain.AggSum
			if !useSum && len(kinds) == 1 && !kinds[domain.TipeRatio] {
				useSum = true
			}
			if useSum {
				if kinds[domain.TipeScalar] && !kinds[domain.TipeCount] {
					c.kind = domain.TipeScalar
				} else {
					c.kind = domain.TipeCount
				}
				c.num, c.den = sum, 0
				if any {
					v := sum
					c.value = &v
				}
			} else {
				c.kind = domain.TipeRatio
				if ratioScale == 0 {
					ratioScale = 100
				}
				c.scale = ratioScale
				c.num, c.den = ratioNum, ratioDen
				if c.den > 0 {
					v := c.num / c.den * ratioScale
					c.value = &v
				}
			}
			switch {
			case allComplete:
				c.status = domain.ValueComplete
			case any:
				c.status = domain.ValuePartial
			default:
				c.status = domain.ValueNoData
			}
		} else {
			c.status = domain.ValueNoData
		}
		memo[labelID] = c
		return c
	}
	hashes, err := s.values.ListFinalFileHashes(ctx, period.ID, scopeType, scopeID)
	if err != nil {
		s.values.FinishCalculationRun(ctx, run.ID, "failed", err.Error(), now())
		return nil, err
	}
	sourceHash := hashHashes(hashes)
	dataVersion := fmt.Sprintf("v%s-%s", period.Kode, run.ID[:8])
	writes := make([]domain.IKUValueWrite, 0, len(labels))
	for _, label := range labels {
		c := compute(label.ID)
		value := domain.IKUValue{
			LabelID: label.ID, PeriodeID: period.ID, ScopeType: scopeType, ScopeID: strPtr(scopeID),
			Numerator: c.num, Denominator: c.den, Value: c.value, Status: c.status,
			FormulaVersion: "v1", CalculationRunID: &run.ID, SourceHash: &sourceHash,
			DataVersion: &dataVersion, UpdatedAt: now(),
		}
		writes = append(writes, domain.IKUValueWrite{Value: value, Details: c.details})
	}
	if err := s.values.UpsertIKUValues(ctx, writes); err != nil {
		s.values.FinishCalculationRun(ctx, run.ID, "failed", err.Error(), now())
		return nil, err
	}
	if err := s.values.FinishCalculationRun(ctx, run.ID, "done", "", now()); err != nil {
		return nil, err
	}
	run.Status = "done"
	return run, nil
}

func (s *CalculationService) Values(ctx context.Context, periodeID, scopeType, scopeID string) ([]domain.IKUValue, error) {
	if scopeType == "" {
		scopeType = "KABUPATEN"
	}
	if scopeID == "" {
		scopeID = "wil-kab-swb"
	}
	return s.values.ListIKUValues(ctx, periodeID, scopeType, scopeID)
}

func (s *CalculationService) ValueByKode(ctx context.Context, kode, periodeID, scopeType, scopeID string) (*domain.IKUValue, error) {
	labels, err := s.master.ListLabels(ctx)
	if err != nil {
		return nil, err
	}
	for _, l := range labels {
		if l.Kode == kode {
			return s.values.GetIKUValue(ctx, l.ID, periodeID, scopeType, scopeID)
		}
	}
	return nil, domain.NotFound("IKU")
}

func (s *CalculationService) Detail(ctx context.Context, ikuValueID string) ([]domain.IKUValueDetail, error) {
	return s.values.ListIKUValueDetails(ctx, ikuValueID)
}

func (s *CalculationService) LatestRun(ctx context.Context, periodeID string) (*domain.CalculationRun, error) {
	return s.values.LatestCalculationRun(ctx, periodeID)
}

func ratioMultiplier(satuan string) float64 {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(satuan)), "rasio") {
		return 1
	}
	return 100
}

func hashHashes(hashes []string) string {
	if len(hashes) == 0 {
		return ""
	}
	sort.Strings(hashes)
	sum := sha256.Sum256([]byte(strings.Join(hashes, "|")))
	return hex.EncodeToString(sum[:])
}
