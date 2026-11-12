package domain

type SubmissionFilter struct {
	Status    string
	PeriodeID string
	FaskesID  string
	Limit     int
}

type AuditFilter struct {
	EntityType string
	EntityID   string
	Limit      int
}

type DashboardFilter struct {
	PeriodeID   string
	ScopeType   string
	ScopeID     string
	FaskesID    string
	BidangID    string
	IndikatorID string
	Status      string
}

type DashboardSPMItem struct {
	Kode             string   `json:"kode"`
	Indikator        string   `json:"indikator"`
	LabelID          *string  `json:"label_id,omitempty"`
	IndikatorID      *string  `json:"indikator_id,omitempty"`
	LabelKode        *string  `json:"label_kode,omitempty"`
	LabelNama        *string  `json:"label_nama,omitempty"`
	IKUValueID       *string  `json:"iku_value_id,omitempty"`
	IndicatorValueID *string  `json:"indicator_value_id,omitempty"`
	CalculationRunID *string  `json:"calculation_run_id,omitempty"`
	SourceHash       *string  `json:"source_hash,omitempty"`
	DataVersion      *string  `json:"data_version,omitempty"`
	Target           float64  `json:"target"`
	BaselinePrev     *float64 `json:"baseline_prev,omitempty"`
	BaselineYear     int      `json:"baseline_year,omitempty"`
	BaselineTW1      *float64 `json:"baseline_tw1,omitempty"`
	Numerator        *float64 `json:"numerator,omitempty"`
	Denominator      *float64 `json:"denominator,omitempty"`
	Value            *float64 `json:"value,omitempty"`
	ValueStatus      string   `json:"value_status"`
	DisplayStatus    string   `json:"display_status"`
	Gap              *float64 `json:"gap,omitempty"`
	UpdatedAt        *string  `json:"updated_at,omitempty"`
}

type SPMTrace struct {
	Row            DashboardSPMItem   `json:"row"`
	IKUValue       *IKUValue          `json:"iku_value,omitempty"`
	Details        []IKUValueDetail   `json:"details,omitempty"`
	IndicatorValue *IndicatorValue    `json:"indicator_value,omitempty"`
	SubmissionRow  *SubmissionRow     `json:"submission_row,omitempty"`
	Submission     *Submission        `json:"submission,omitempty"`
	File           *FileObject        `json:"file,omitempty"`
	Reviews        []SubmissionReview `json:"reviews,omitempty"`
}

type CreateUploadParams struct {
	SubmissionID      string
	Kode              string
	PeriodKode        string
	UploaderID        string
	FaskesID          string
	PeriodeID         string
	MasterReleaseID   string
	WorkflowVersionID string
	File              FileObject
	IdempotencyKey    string
	Now               string
}

type SubmissionRowInput struct {
	ID                 string
	RowNo              int
	KodeIKURaw         string
	KodeIndikatorRaw   string
	MatchedIndikatorID string
	MatchedLabelID     string
	ScopeType          string
	ScopeCode          string
	Value              float64
	Numerator          float64
	Denominator        float64
	Classification     ClassificationStatus
	ErrorMessage       string
}

type IndicatorValueInput struct {
	SubmissionRowID string
	IndikatorID     string
	PeriodeID       string
	FaskesID        string
	ScopeType       string
	ScopeID         string
	Numerator       float64
	Denominator     float64
	Value           float64
	Status          string
	DataVersion     string
}

type IKUValueWrite struct {
	Value   IKUValue
	Details []IKUValueDetail
}

type DashboardSummary struct {
	Total        int      `json:"total"`
	OnTarget     int      `json:"on_target"`
	BelowTarget  int      `json:"below_target"`
	NoData       int      `json:"no_data"`
	Average      float64  `json:"average"`
	ParentValue  *float64 `json:"parent_value,omitempty"`
	ParentStatus string   `json:"parent_status"`
}

type SPMFaskesDetail struct {
	FaskesID       string   `json:"faskes_id"`
	FaskesKode     string   `json:"faskes_kode"`
	FaskesNama     string   `json:"faskes_nama"`
	Uploaded       bool     `json:"uploaded"`
	SubmissionID   *string  `json:"submission_id,omitempty"`
	SubmissionKode *string  `json:"submission_kode,omitempty"`
	FileID         *string  `json:"file_id,omitempty"`
	FileName       *string  `json:"file_name,omitempty"`
	UploaderNama   *string  `json:"uploader_nama,omitempty"`
	Numerator      *float64 `json:"numerator,omitempty"`
	Denominator    *float64 `json:"denominator,omitempty"`
	Value          *float64 `json:"value,omitempty"`
	Status         string   `json:"status"`
	UpdatedAt      *string  `json:"updated_at,omitempty"`
}

type SPMDetail struct {
	Row        DashboardSPMItem  `json:"row"`
	Faskes     []SPMFaskesDetail `json:"faskes"`
	Provenance []IKUValueDetail  `json:"provenance,omitempty"`
}

type ChartFaskes struct {
	ID   string `json:"id"`
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type ChartSeries struct {
	Kode   string     `json:"kode"`
	Nama   string     `json:"nama"`
	Values []*float64 `json:"values"`
}

type IKUChart struct {
	Faskes []ChartFaskes `json:"faskes"`
	Series []ChartSeries `json:"series"`
}
