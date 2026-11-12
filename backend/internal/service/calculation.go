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
	scopes := []scope{{Type: "KABUPATEN", ID: kabupatenID}}
	allValues, err := s.values.ListFinalIndicatorValues(ctx, periodeID)
	if err != nil {
		return nil, err
	}
	seenFaskes := map[string]bool{}
	for _, v := range allValues {
		if v.FaskesID != nil && *v.FaskesID != "" && !seenFaskes[*v.FaskesID] {
			seenFaskes[*v.FaskesID] = true
			scopes = append(scopes, scope{Type: "FASKES", ID: *v.FaskesID})
		}
	}
	var runs []domain.CalculationRun
	for _, sc := range scopes {
		run, err := s.calculateScope(ctx, period, sc.Type, sc.ID, trigger)
		if err != nil {
			return runs, err
		}
		runs = append(runs, *run)
	}
	return runs, nil
}

func (s *CalculationService) calculateScope(ctx context.Context, period *domain.Periode, scopeType, scopeID, trigger string) (*domain.CalculationRun, error) {
	run := &domain.CalculationRun{
		PeriodeID: period.ID, ScopeType: scopeType, ScopeID: strPtr(scopeID),
		Trigger: trigger, Status: "running", StartedAt: now(),
	}
	if err := s.values.CreateCalculationRun(ctx, run); err != nil {
		return nil, err
	}
	values, err := s.values.ListFinalIndicatorValuesForScope(ctx, period.ID, scopeType, scopeID)
	if err != nil {
		s.values.FinishCalculationRun(ctx, run.ID, "failed", err.Error(), now())
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
			for _, indicatorID := range indicatorLinks[labelID] {
				for _, iv := range indicatorValues[indicatorID] {
					c.num += iv.Numerator
					c.den += iv.Denominator
					v := iv.Numerator
					if iv.Denominator != 0 {
						v = iv.Numerator / iv.Denominator * 100
					}
					c.details = append(c.details, domain.IKUValueDetail{
						ComponentType: "INDIKATOR", ComponentID: iv.IndikatorID,
						ComponentValue: &v, Numerator: &iv.Numerator, Denominator: &iv.Denominator,
						Weight: 1, Contribution: &v,
					})
				}
			}
			if c.den > 0 {
				c.status = domain.ValueComplete
			} else {
				c.status = domain.ValueNoData
			}
		} else if len(children[labelID]) > 0 {
			allComplete := true
			any := false
			for _, comp := range children[labelID] {
				childID := comp.KomponenID
				child := compute(childID)
				if child.status == domain.ValueComplete {
					any = true
				} else {
					allComplete = false
				}
				if child.num > 0 || child.den > 0 {
					c.num += child.num
					c.den += child.den
				}
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
		if c.den > 0 {
			v := c.num / c.den * 100
			c.value = &v
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

func hashHashes(hashes []string) string {
	if len(hashes) == 0 {
		return ""
	}
	sort.Strings(hashes)
	sum := sha256.Sum256([]byte(strings.Join(hashes, "|")))
	return hex.EncodeToString(sum[:])
}
