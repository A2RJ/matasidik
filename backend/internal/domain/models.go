package domain

// Semua ID memakai UUID string agar mudah dipindahkan ke PostgreSQL.
// Waktu disimpan sebagai RFC3339 dalam TEXT (SQLite) dan tetap kompatibel
// dengan timestamptz saat migrasi dilakukan.

type User struct {
	ID           string  `db:"id" json:"id"`
	Nama         string  `db:"nama" json:"nama"`
	Username     string  `db:"username" json:"username"`
	Role         Role    `db:"role" json:"role"`
	FaskesID     *string `db:"faskes_id" json:"faskes_id,omitempty"`
	BidangID     *string `db:"bidang_id" json:"bidang_id,omitempty"`
	PositionID   *string `db:"position_id" json:"position_id,omitempty"`
	PositionNama *string `db:"position_nama" json:"position_nama,omitempty"`
	LastLoginAt  *string `db:"last_login_at" json:"last_login_at,omitempty"`
	Aktif        bool    `db:"aktif" json:"aktif"`
	CreatedAt    string  `db:"created_at" json:"created_at"`
}

type UserAuth struct {
	User
	PasswordHash string `db:"password_hash" json:"-"`
}

type Position struct {
	ID        string  `db:"id" json:"id"`
	Kode      string  `db:"kode" json:"kode"`
	Nama      string  `db:"nama" json:"nama"`
	Role      Role    `db:"role" json:"role"`
	Level     int     `db:"level" json:"level"`
	BidangID  *string `db:"bidang_id" json:"bidang_id,omitempty"`
	FaskesID  *string `db:"faskes_id" json:"faskes_id,omitempty"`
	Aktif     bool    `db:"aktif" json:"aktif"`
	CreatedAt string  `db:"created_at" json:"created_at"`
}

type AuthSession struct {
	ID          string  `db:"id" json:"id"`
	UserID      string  `db:"user_id" json:"user_id"`
	RefreshHash string  `db:"refresh_hash" json:"-"`
	UserAgent   *string `db:"user_agent" json:"user_agent,omitempty"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
	ExpiresAt   string  `db:"expires_at" json:"expires_at"`
	RevokedAt   *string `db:"revoked_at" json:"revoked_at,omitempty"`
}

type Wilayah struct {
	ID       string  `db:"id" json:"id"`
	Kode     string  `db:"kode" json:"kode"`
	Nama     string  `db:"nama" json:"nama"`
	Level    string  `db:"level" json:"level"`
	ParentID *string `db:"parent_id" json:"parent_id,omitempty"`
}

type Faskes struct {
	ID        string  `db:"id" json:"id"`
	Kode      string  `db:"kode" json:"kode"`
	Nama      string  `db:"nama" json:"nama"`
	WilayahID *string `db:"wilayah_id" json:"wilayah_id,omitempty"`
	Jenis     string  `db:"jenis" json:"jenis"`
	Aktif     bool    `db:"aktif" json:"aktif"`
}

type Bidang struct {
	ID       string  `db:"id" json:"id"`
	Kode     string  `db:"kode" json:"kode"`
	Nama     string  `db:"nama" json:"nama"`
	Jenis    string  `db:"jenis" json:"jenis"`
	ParentID *string `db:"parent_id" json:"parent_id,omitempty"`
	Aktif    bool    `db:"aktif" json:"aktif"`
}

type DashboardCategory struct {
	ID         string  `db:"id" json:"id"`
	Kode       string  `db:"kode" json:"kode"`
	Nama       string  `db:"nama" json:"nama"`
	Jenis      string  `db:"jenis" json:"jenis"`
	Urutan     int     `db:"urutan" json:"urutan"`
	Aktif      bool    `db:"aktif" json:"aktif"`
	IsSystem   bool    `db:"is_system" json:"is_system"`
	ConfigJSON *string `db:"config_json" json:"config_json,omitempty"`
	CreatedAt  string  `db:"created_at" json:"created_at"`
	UpdatedAt  string  `db:"updated_at" json:"updated_at"`
}

type Periode struct {
	ID            string `db:"id" json:"id"`
	Tahun         int    `db:"tahun" json:"tahun"`
	Bulan         int    `db:"bulan" json:"bulan"`
	Kode          string `db:"kode" json:"kode"`
	Status        string `db:"status" json:"status"`
	IsLateAllowed bool   `db:"is_late_allowed" json:"is_late_allowed"`
}

type MasterRelease struct {
	ID          string  `db:"id" json:"id"`
	Tahun       int     `db:"tahun" json:"tahun"`
	Status      string  `db:"status" json:"status"`
	PublishedAt *string `db:"published_at" json:"published_at,omitempty"`
	PublishedBy *string `db:"published_by" json:"published_by,omitempty"`
}

type Label struct {
	ID                string            `db:"id" json:"id"`
	Kode              string            `db:"kode" json:"kode"`
	Nama              string            `db:"nama" json:"nama"`
	Tipe              string            `db:"tipe" json:"tipe"`
	CategoryID        *string           `db:"category_id" json:"category_id,omitempty"`
	ParentID          *string           `db:"parent_id" json:"parent_id,omitempty"`
	MasterReleaseID   string            `db:"master_release_id" json:"master_release_id"`
	AggregationMethod AggregationMethod `db:"aggregation_method" json:"aggregation_method"`
	AllowDirectInput  bool              `db:"allow_direct_input" json:"allow_direct_input"`
	IsRollupOnly      bool              `db:"is_rollup_only" json:"is_rollup_only"`
	PJBidangID        *string           `db:"pj_bidang_id" json:"pj_bidang_id,omitempty"`
	IKUOwnerID        *string           `db:"iku_owner_id" json:"iku_owner_id,omitempty"`
	Level             int               `db:"level" json:"level"`
	Path              *string           `db:"path" json:"path,omitempty"`
	Aktif             bool              `db:"aktif" json:"aktif"`
	BerlakuDari       *string           `db:"berlaku_dari" json:"berlaku_dari,omitempty"`
	BerlakuSampai     *string           `db:"berlaku_sampai" json:"berlaku_sampai,omitempty"`
	// Join-only fields.
	CategoryKode *string `db:"category_kode" json:"category_kode,omitempty"`
	CategoryNama *string `db:"category_nama" json:"category_nama,omitempty"`
}

type Indikator struct {
	ID               string  `db:"id" json:"id"`
	Kode             string  `db:"kode" json:"kode"`
	Nama             string  `db:"nama" json:"nama"`
	Satuan           string  `db:"satuan" json:"satuan"`
	Definisi         *string `db:"definisi" json:"definisi,omitempty"`
	NumeratorLabel   *string `db:"numerator_label" json:"numerator_label,omitempty"`
	DenominatorLabel *string `db:"denominator_label" json:"denominator_label,omitempty"`
	SifatInput       string  `db:"sifat_input" json:"sifat_input"`
	Frekuensi        string  `db:"frekuensi" json:"frekuensi"`
	PJBidangID       *string `db:"pj_bidang_id" json:"pj_bidang_id,omitempty"`
	PJTimkerID       *string `db:"pj_timker_id" json:"pj_timker_id,omitempty"`
	PJTimker         *string `db:"pj_timker" json:"pj_timker,omitempty"`
	Aktif            bool    `db:"aktif" json:"aktif"`
	// Join-only fields for master preview.
	LabelIDs       []string `db:"-" json:"label_ids,omitempty"`
	PrimaryLabelID *string  `db:"-" json:"primary_label_id,omitempty"`
}

type IndikatorLabel struct {
	IndikatorID string  `db:"indikator_id" json:"indikator_id"`
	LabelID     string  `db:"label_id" json:"label_id"`
	Bobot       float64 `db:"bobot" json:"bobot"`
	IsPrimary   bool    `db:"is_primary" json:"is_primary"`
}

type IKUComponent struct {
	ID            string  `db:"id" json:"id"`
	ParentLabelID string  `db:"parent_label_id" json:"parent_label_id"`
	TipeKomponen  string  `db:"tipe_komponen" json:"tipe_komponen"`
	KomponenID    string  `db:"komponen_id" json:"komponen_id"`
	Bobot         float64 `db:"bobot" json:"bobot"`
	Sort          int     `db:"sort" json:"sort"`
}

type Target struct {
	ID                   string   `db:"id" json:"id"`
	LabelID              *string  `db:"label_id" json:"label_id,omitempty"`
	IndikatorID          *string  `db:"indikator_id" json:"indikator_id,omitempty"`
	PeriodeID            string   `db:"periode_id" json:"periode_id"`
	ScopeType            string   `db:"scope_type" json:"scope_type"`
	ScopeID              *string  `db:"scope_id" json:"scope_id,omitempty"`
	TargetTahunan        float64  `db:"target_tahunan" json:"target_tahunan"`
	TargetBulanan        float64  `db:"target_bulanan" json:"target_bulanan"`
	TargetLayananTahunan *float64 `db:"target_layanan_tahunan" json:"target_layanan_tahunan,omitempty"`
	TargetLayananBulanan *float64 `db:"target_layanan_bulanan" json:"target_layanan_bulanan,omitempty"`
}

type CapaianTahunan struct {
	ID            string  `db:"id" json:"id"`
	LabelID       *string `db:"label_id" json:"label_id,omitempty"`
	IndikatorID   *string `db:"indikator_id" json:"indikator_id,omitempty"`
	Tahun         int     `db:"tahun" json:"tahun"`
	ScopeType     string  `db:"scope_type" json:"scope_type"`
	ScopeID       *string `db:"scope_id" json:"scope_id,omitempty"`
	Nilai         float64 `db:"nilai" json:"nilai"`
	Catatan       *string `db:"catatan" json:"catatan,omitempty"`
	UpdatedBy     *string `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt     string  `db:"created_at" json:"created_at"`
	UpdatedAt     string  `db:"updated_at" json:"updated_at"`
	LabelKode     *string `db:"label_kode" json:"label_kode,omitempty"`
	LabelNama     *string `db:"label_nama" json:"label_nama,omitempty"`
	IndikatorKode *string `db:"indikator_kode" json:"indikator_kode,omitempty"`
	IndikatorNama *string `db:"indikator_nama" json:"indikator_nama,omitempty"`
	UpdatedByNama *string `db:"updated_by_nama" json:"updated_by_nama,omitempty"`
}

type WorkflowVersion struct {
	ID          string  `db:"id" json:"id"`
	Kode        string  `db:"kode" json:"kode"`
	VersionNo   int     `db:"version_no" json:"version_no"`
	StepsJSON   string  `db:"steps_json" json:"steps_json"`
	Status      string  `db:"status" json:"status"`
	PublishedAt *string `db:"published_at" json:"published_at,omitempty"`
}

type FileObject struct {
	ID           string  `db:"id" json:"id"`
	OriginalName string  `db:"original_name" json:"original_name"`
	SHA256       string  `db:"sha256" json:"sha256"`
	Path         string  `db:"path" json:"path"`
	Mime         *string `db:"mime" json:"mime,omitempty"`
	Size         int64   `db:"size" json:"size"`
	UploadedBy   *string `db:"uploaded_by" json:"uploaded_by,omitempty"`
	ScanStatus   string  `db:"scan_status" json:"scan_status"`
	CreatedAt    string  `db:"created_at" json:"created_at"`
}

type Submission struct {
	ID                string           `db:"id" json:"id"`
	Kode              string           `db:"kode" json:"kode"`
	Tipe              string           `db:"tipe" json:"tipe"`
	UploaderID        string           `db:"uploader_id" json:"uploader_id"`
	FaskesID          string           `db:"faskes_id" json:"faskes_id"`
	PeriodeID         string           `db:"periode_id" json:"periode_id"`
	FileID            string           `db:"file_id" json:"file_id"`
	FileHash          string           `db:"file_hash" json:"file_hash"`
	Status            SubmissionStatus `db:"status" json:"status"`
	Step              WorkflowStep     `db:"step" json:"step"`
	MasterReleaseID   *string          `db:"master_release_id" json:"master_release_id,omitempty"`
	WorkflowVersionID *string          `db:"workflow_version_id" json:"workflow_version_id,omitempty"`
	IdempotencyKey    *string          `db:"idempotency_key" json:"idempotency_key,omitempty"`
	CreatedAt         string           `db:"created_at" json:"created_at"`
	UpdatedAt         string           `db:"updated_at" json:"updated_at"`
	// Join-only fields.
	UploaderNama *string `db:"uploader_nama" json:"uploader_nama,omitempty"`
	FaskesNama   *string `db:"faskes_nama" json:"faskes_nama,omitempty"`
	PeriodeKode  *string `db:"periode_kode" json:"periode_kode,omitempty"`
}

type SubmissionRevision struct {
	ID               string  `db:"id" json:"id"`
	SubmissionID     string  `db:"submission_id" json:"submission_id"`
	RevisionNo       int     `db:"revision_no" json:"revision_no"`
	ParentRevisionID *string `db:"parent_revision_id" json:"parent_revision_id,omitempty"`
	FileID           string  `db:"file_id" json:"file_id"`
	CreatedBy        *string `db:"created_by" json:"created_by,omitempty"`
	CreatedAt        string  `db:"created_at" json:"created_at"`
}

type SubmissionGroup struct {
	ID             string  `db:"id" json:"id"`
	SubmissionID   string  `db:"submission_id" json:"submission_id"`
	OwnerBidangID  *string `db:"owner_bidang_id" json:"owner_bidang_id,omitempty"`
	PrimaryLabelID *string `db:"primary_label_id" json:"primary_label_id,omitempty"`
	Status         string  `db:"status" json:"status"`
	CurrentStep    string  `db:"current_step" json:"current_step"`
}

type SubmissionRow struct {
	ID                 string               `db:"id" json:"id"`
	SubmissionID       string               `db:"submission_id" json:"submission_id"`
	GroupID            *string              `db:"group_id" json:"group_id,omitempty"`
	RowNo              int                  `db:"row_no" json:"row_no"`
	KodeIKURaw         *string              `db:"kode_iku_raw" json:"kode_iku_raw,omitempty"`
	KodeIndikatorRaw   *string              `db:"kode_indikator_raw" json:"kode_indikator_raw,omitempty"`
	MatchedIndikatorID *string              `db:"matched_indikator_id" json:"matched_indikator_id,omitempty"`
	MatchedLabelID     *string              `db:"matched_label_id" json:"matched_label_id,omitempty"`
	ScopeType          string               `db:"scope_type" json:"scope_type"`
	ScopeCode          *string              `db:"scope_code" json:"scope_code,omitempty"`
	Value              *float64             `db:"nilai" json:"nilai,omitempty"`
	Numerator          *float64             `db:"numerator" json:"numerator,omitempty"`
	Denominator        *float64             `db:"denominator" json:"denominator,omitempty"`
	Classification     ClassificationStatus `db:"classification_status" json:"classification_status"`
	ErrorMessage       *string              `db:"error_message" json:"error_message,omitempty"`
	CreatedAt          string               `db:"created_at" json:"created_at"`
	// Join-only fields.
	IndikatorKode *string `db:"indikator_kode" json:"indikator_kode,omitempty"`
	IndikatorNama *string `db:"indikator_nama" json:"indikator_nama,omitempty"`
	IKUKode       *string `db:"iku_kode" json:"iku_kode,omitempty"`
	IKUNama       *string `db:"iku_nama" json:"iku_nama,omitempty"`
}

type SubmissionReview struct {
	ID           string  `db:"id" json:"id"`
	SubmissionID string  `db:"submission_id" json:"submission_id"`
	GroupID      *string `db:"group_id" json:"group_id,omitempty"`
	Step         string  `db:"step" json:"step"`
	ReviewerID   *string `db:"reviewer_id" json:"reviewer_id,omitempty"`
	Action       string  `db:"action" json:"action"`
	Notes        *string `db:"notes" json:"notes,omitempty"`
	DetailsJSON  *string `db:"details_json" json:"details_json,omitempty"`
	CreatedAt    string  `db:"created_at" json:"created_at"`
	ReviewerNama *string `db:"reviewer_nama" json:"reviewer_nama,omitempty"`
}

type Job struct {
	ID           string    `db:"id" json:"id"`
	Type         string    `db:"type" json:"type"`
	PayloadJSON  string    `db:"payload_json" json:"payload_json"`
	DedupeKey    *string   `db:"dedupe_key" json:"dedupe_key,omitempty"`
	Status       JobStatus `db:"status" json:"status"`
	Attempts     int       `db:"attempts" json:"attempts"`
	MaxAttempts  int       `db:"max_attempts" json:"max_attempts"`
	RunAfter     string    `db:"run_after" json:"run_after"`
	LockedBy     *string   `db:"locked_by" json:"locked_by,omitempty"`
	LockedAt     *string   `db:"locked_at" json:"locked_at,omitempty"`
	HeartbeatAt  *string   `db:"heartbeat_at" json:"heartbeat_at,omitempty"`
	LastError    *string   `db:"last_error" json:"last_error,omitempty"`
	ProgressJSON *string   `db:"progress_json" json:"progress_json,omitempty"`
	CreatedAt    string    `db:"created_at" json:"created_at"`
	FinishedAt   *string   `db:"finished_at" json:"finished_at,omitempty"`
}

type IndicatorValue struct {
	ID              string  `db:"id" json:"id"`
	SubmissionRowID string  `db:"submission_row_id" json:"submission_row_id"`
	IndikatorID     string  `db:"indikator_id" json:"indikator_id"`
	PeriodeID       string  `db:"periode_id" json:"periode_id"`
	FaskesID        *string `db:"faskes_id" json:"faskes_id,omitempty"`
	ScopeType       string  `db:"scope_type" json:"scope_type"`
	ScopeID         *string `db:"scope_id" json:"scope_id,omitempty"`
	Numerator       float64 `db:"numerator" json:"numerator"`
	Denominator     float64 `db:"denominator" json:"denominator"`
	Value           float64 `db:"value" json:"value"`
	Status          string  `db:"status" json:"status"`
	DataVersion     *string `db:"data_version" json:"data_version,omitempty"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
}

type CalculationRun struct {
	ID         string  `db:"id" json:"id"`
	PeriodeID  string  `db:"periode_id" json:"periode_id"`
	ScopeType  string  `db:"scope_type" json:"scope_type"`
	ScopeID    *string `db:"scope_id" json:"scope_id,omitempty"`
	Trigger    string  `db:"trigger" json:"trigger"`
	Status     string  `db:"status" json:"status"`
	StartedAt  string  `db:"started_at" json:"started_at"`
	FinishedAt *string `db:"finished_at" json:"finished_at,omitempty"`
	Error      *string `db:"error" json:"error,omitempty"`
}

type IKUValue struct {
	ID               string      `db:"id" json:"id"`
	LabelID          string      `db:"label_id" json:"label_id"`
	PeriodeID        string      `db:"periode_id" json:"periode_id"`
	ScopeType        string      `db:"scope_type" json:"scope_type"`
	ScopeID          *string     `db:"scope_id" json:"scope_id,omitempty"`
	Numerator        float64     `db:"numerator" json:"numerator"`
	Denominator      float64     `db:"denominator" json:"denominator"`
	Value            *float64    `db:"value" json:"value,omitempty"`
	Status           ValueStatus `db:"status" json:"status"`
	FormulaVersion   string      `db:"formula_version" json:"formula_version"`
	CalculationRunID *string     `db:"calculation_run_id" json:"calculation_run_id,omitempty"`
	SourceHash       *string     `db:"source_hash" json:"source_hash,omitempty"`
	DataVersion      *string     `db:"data_version" json:"data_version,omitempty"`
	UpdatedAt        string      `db:"updated_at" json:"updated_at"`
	// Join-only fields.
	LabelKode *string `db:"label_kode" json:"label_kode,omitempty"`
	LabelNama *string `db:"label_nama" json:"label_nama,omitempty"`
	LabelTipe *string `db:"label_tipe" json:"label_tipe,omitempty"`
}

type IKUValueDetail struct {
	ID               string   `db:"id" json:"id"`
	IKUValueID       string   `db:"iku_value_id" json:"iku_value_id"`
	ComponentType    string   `db:"component_type" json:"component_type"`
	ComponentID      string   `db:"component_id" json:"component_id"`
	ComponentLabelID *string  `db:"component_label_id" json:"component_label_id,omitempty"`
	ComponentValue   *float64 `db:"component_value" json:"component_value,omitempty"`
	Numerator        *float64 `db:"numerator" json:"numerator,omitempty"`
	Denominator      *float64 `db:"denominator" json:"denominator,omitempty"`
	Weight           float64  `db:"weight" json:"weight"`
	Contribution     *float64 `db:"contribution" json:"contribution,omitempty"`
	// Join-only fields.
	LabelKode     *string `db:"label_kode" json:"label_kode,omitempty"`
	LabelNama     *string `db:"label_nama" json:"label_nama,omitempty"`
	IndikatorKode *string `db:"indikator_kode" json:"indikator_kode,omitempty"`
	IndikatorNama *string `db:"indikator_nama" json:"indikator_nama,omitempty"`
}

type DashboardRow struct {
	ID           string  `db:"id" json:"id"`
	Kode         string  `db:"kode" json:"kode"`
	LabelID      *string `db:"label_id" json:"label_id,omitempty"`
	IndikatorID  *string `db:"indikator_id" json:"indikator_id,omitempty"`
	Urutan       int     `db:"urutan" json:"urutan"`
	Target       float64 `db:"target" json:"target"`
	BaselineJSON *string `db:"baseline_json" json:"baseline_json,omitempty"`
	FormatJSON   *string `db:"format_json" json:"format_json,omitempty"`
}

type Faktor struct {
	ID         string  `db:"id" json:"id"`
	IKUValueID string  `db:"iku_value_id" json:"iku_value_id"`
	PeriodeID  string  `db:"periode_id" json:"periode_id"`
	Pendukung  *string `db:"pendukung" json:"pendukung,omitempty"`
	Penghambat *string `db:"penghambat" json:"penghambat,omitempty"`
	CreatedBy  *string `db:"created_by" json:"created_by,omitempty"`
	UpdatedAt  string  `db:"updated_at" json:"updated_at"`
	// Join-only fields.
	LabelKode *string `db:"label_kode" json:"label_kode,omitempty"`
	LabelNama *string `db:"label_nama" json:"label_nama,omitempty"`
}

type TindakLanjut struct {
	ID          string   `db:"id" json:"id"`
	IKUValueID  *string  `db:"iku_value_id" json:"iku_value_id,omitempty"`
	PeriodeID   string   `db:"periode_id" json:"periode_id"`
	Masalah     string   `db:"masalah" json:"masalah"`
	Rekomendasi string   `db:"rekomendasi" json:"rekomendasi"`
	PJUserID    *string  `db:"pj_user_id" json:"pj_user_id,omitempty"`
	Status      TLStatus `db:"status" json:"status"`
	CreatedBy   *string  `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   string   `db:"created_at" json:"created_at"`
	UpdatedAt   string   `db:"updated_at" json:"updated_at"`
	// Join-only fields.
	PJNama    *string `db:"pj_nama" json:"pj_nama,omitempty"`
	LabelKode *string `db:"label_kode" json:"label_kode,omitempty"`
	LabelNama *string `db:"label_nama" json:"label_nama,omitempty"`
}

type AuditLog struct {
	ID               string  `db:"id" json:"id"`
	ActorID          *string `db:"actor_id" json:"actor_id,omitempty"`
	Action           string  `db:"action" json:"action"`
	EntityType       string  `db:"entity_type" json:"entity_type"`
	EntityID         *string `db:"entity_id" json:"entity_id,omitempty"`
	BeforeJSON       *string `db:"before_json" json:"before_json,omitempty"`
	AfterJSON        *string `db:"after_json" json:"after_json,omitempty"`
	Reason           *string `db:"reason" json:"reason,omitempty"`
	RequestID        *string `db:"request_id" json:"request_id,omitempty"`
	ImpersonatorID   *string `db:"impersonator_id" json:"impersonator_id,omitempty"`
	CreatedAt        string  `db:"created_at" json:"created_at"`
	ActorNama        *string `db:"actor_nama" json:"actor_nama,omitempty"`
	ImpersonatorNama *string `db:"impersonator_nama" json:"impersonator_nama,omitempty"`
}
